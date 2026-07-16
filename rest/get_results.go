package rest

import (
	"iter"
	"math"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/boggydigital/redux"
	"github.com/boggydigital/strom"
	"github.com/boggydigital/strom/vars"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/rest/view_models"
	"github.com/boggydigital/yet/yeti"
	"github.com/boggydigital/yet_urls/youtube_urls"
)

type ResultsViewModel struct {
	SearchQuery string
	Refinements []string
	Channels    []*view_models.ChannelViewModel
	Playlists   []*view_models.PlaylistViewModel
	Videos      []*view_models.VideoViewModel
}

var propertyTitles = map[string]string{
	data.VideoOwnerChannelNameProperty:  "Channel",
	data.VideoEndedDateProperty:         "Ended",
	data.VideoPublishDateProperty:       "Published",
	data.VideoDownloadCompletedProperty: "Downloaded",
}

var propertiesOrder = []string{
	data.VideoOwnerChannelNameProperty,
	data.VideoEndedDateProperty,
	data.VideoPublishDateProperty,
	data.VideoDownloadCompletedProperty,
}

func GetResults(w http.ResponseWriter, r *http.Request) {

	var err error
	rdx, err = rdx.RefreshWriter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	searchQuery := r.URL.Query().Get("search-query")

	root, body := strom.RootBody("Search")

	body.AddClass("d-f", "fd-c", "rg-n")

	body.Append(navButton("Home", "/"))

	body.Append(strom.CreateText("h1", "Results for '"+searchQuery+"'"))

	sid, err := youtube_urls.GetSearchResultsPage(
		http.DefaultClient,
		strings.Split(searchQuery, " ")...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	propertyValues := extractSearchVideosMetadata(sid.VideoRenderers())
	for property, keyValues := range propertyValues {
		if err = rdx.BatchAddValues(property, keyValues); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	propertyValues = extractSearchPlaylistMetadata(sid.PlaylistRenderers())
	for property, keyValues := range propertyValues {
		if err = rdx.BatchAddValues(property, keyValues); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	propertyValues = extractSearchChannelMetadata(sid.ChannelRenderers())
	for property, keyValues := range propertyValues {
		if err = rdx.BatchAddValues(property, keyValues); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	hasRefinements := len(sid.Refinements) > 0
	hasChannels := len(sid.ChannelRenderers()) > 0
	hasPlaylists := len(sid.PlaylistRenderers()) > 0
	hasVideos := len(sid.VideoRenderers()) > 0

	if hasRefinements {
		body.Append(strom.CreateText("h2", "Refinements"))
	}

	if hasChannels {
		body.Append(strom.CreateText("h2", "Channels"))

		channels := strom.Create("ul", "d-f", "cg-n", "rg-n").
			SetStyle(map[string]string{
				"flex-flow": "row wrap",
			})
		body.Append(channels)

		for _, chr := range sid.ChannelRenderers() {
			channels.Append(channelTile(chr.ChannelId, rdx))
		}
	}

	if hasPlaylists {
		//
	}

	if hasVideos {
		if hasRefinements || hasChannels || hasPlaylists {
			body.Append(strom.CreateText("h2", "Videos"))
		}

		videos := strom.Create("ul", "d-f", "cg-n", "rg-n").
			SetStyle(map[string]string{
				"flex-flow": "row wrap",
			})
		body.Append(videos)

		srv := new(searchResultsVideos{searchInitialData: sid, rdx: rdx})

		videos.Append(strom.Defer(srv.getVideos))
	}

	if err = strom.WriteResponse(w, root); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type searchResultsVideos struct {
	searchInitialData *youtube_urls.SearchInitialData
	rdx               redux.Readable
}

func (vd *searchResultsVideos) getVideos() iter.Seq[strom.Element] {
	return func(yield func(strom.Element) bool) {
		for _, vr := range vd.searchInitialData.VideoRenderers() {
			if !yield(videoTile(vr.VideoId, rdx)) {
				return
			}
		}
	}
}

func videoTile(videoId string, rdx redux.Readable) strom.Element {

	tileContainer := strom.Create("a", "d-f", "fd-c", "rg-n").
		SetAttribute("href", path.Join("/watch", videoId)).
		SetStyle(map[string]string{
			"width":    "calc(1.5 * " + vars.Size(vars.SizeXXXLarge) + ")",
			"position": "relative",
		})

	var ended bool
	if rdx.HasKey(data.VideoEndedDateProperty, videoId) {
		ended = true
	}

	poster := strom.Create("img", "br-s").
		SetAttribute("src", path.Join("/poster?v="+videoId+"&q=hqdefault")).
		SetAttribute("loading", "lazy").
		SetStyle(map[string]string{
			"aspect-ratio": "16/9",
			"width":        "100%",
			"object-fit":   "cover",
		})

	tileContainer.Append(poster)

	if ended {
		poster.SetStyle(map[string]string{
			"filter": "grayscale(1.0)",
		})

		reason := data.Completed
		if ver, ok := rdx.GetLastVal(data.VideoEndedReasonProperty, videoId); ok && ver != "" {
			reason = data.ParseVideoEndedReason(ver)
		}

		tileContainer.Append(strom.CreateText("span", reason.String()).
			SetStyle(map[string]string{
				"position":                  "absolute",
				"top":                       "0",
				"right":                     "0",
				"font-size":                 vars.FontSize(vars.SizeXSmall),
				"padding":                   vars.Size(vars.SizeSmall),
				"border-bottom-left-radius": vars.Size(vars.SizeSmall),
				"border-top-right-radius":   vars.Size(vars.SizeSmall),
				"background-color":          vars.Color(vars.ColorBackground),
			}))
	}

	if durs, sure := rdx.GetLastVal(data.VideoDurationProperty, videoId); sure && durs != "" && durs != "0" {
		if duri, err := strconv.ParseInt(durs, 10, 64); err == nil {

			var remaining int64

			if cts, ok := rdx.GetLastVal(data.VideoProgressProperty, videoId); ok && cts != "" {
				var cti int64
				if cti, err = strconv.ParseInt(cts, 10, 64); err == nil {
					remaining = duri - cti
				}
			}

			durationItems := strom.Create("span", "d-f", "fd-r", "cg-s", "fs-s").
				SetStyle(map[string]string{
					"position":                   "absolute",
					"top":                        "0",
					"left":                       "0",
					"font-size":                  vars.FontSize(vars.SizeXSmall),
					"padding":                    vars.Size(vars.SizeSmall),
					"border-bottom-right-radius": vars.Size(vars.SizeSmall),
					"border-top-left-radius":     vars.Size(vars.SizeSmall),
					"background-color":           vars.Color(vars.ColorBackground),
				})

			durSpan := strom.CreateText("span", formatSeconds(duri)).
				SetStyle(map[string]string{
					"font-size": vars.FontSize(vars.SizeXSmall),
				})

			if remaining > 0 {
				remSpan := strom.CreateText("span", formatSeconds(remaining), "fw-b")
				durationItems.Append(remSpan)
				durSpan.SetStyle(map[string]string{"color": vars.Color(vars.ColorGray)})
			} else {
				if !ended {
					durSpan.AddClass("fw-b")
				}
			}

			durationItems.Append(durSpan)

			tileContainer.Append(durationItems)
		}
	}

	titlePropertiesStack := strom.Create("ul", "d-f", "fd-c", "rg-s")

	if title, ok := rdx.GetLastVal(data.VideoTitleProperty, videoId); ok && title != "" {
		titlePropertiesStack.Append(strom.CreateText("h3", title))
	}

	vsp := videoSummaryProperties(videoId, rdx)

	propertiesStack := strom.Create("ul", "d-f", "fd-c", "rg-xs")
	titlePropertiesStack.Append(propertiesStack)

	for _, p := range propertiesOrder {
		v := vsp[p]
		if v == "" {
			continue
		}

		ptv := propertyTitles[p] + ": " + v

		propertyRow := strom.CreateText("span", ptv).
			SetStyle(map[string]string{
				"color":     vars.Color(vars.ColorGray),
				"font-size": vars.FontSize(vars.SizeXSmall),
			})

		propertiesStack.Append(propertyRow)
	}

	tileContainer.Append(titlePropertiesStack)

	return tileContainer
}

func videoSummaryProperties(videoId string, rdx redux.Readable) map[string]string {
	properties := make(map[string]string)

	if och, ok := rdx.GetLastVal(data.VideoOwnerChannelNameProperty, videoId); ok && och != "" {
		properties[data.VideoOwnerChannelNameProperty] = och
	}

	if ets, ok := rdx.GetLastVal(data.VideoEndedDateProperty, videoId); ok && ets != "" {
		properties[data.VideoEndedDateProperty] = parseAndFormatDate(ets)
	}

	if len(properties) < 2 {
		var publishedDate string
		if pds, ok := rdx.GetLastVal(data.VideoPublishDateProperty, videoId); ok && pds != "" {
			publishedDate = parseAndFormatDate(pds)
		} else {
			if ptts, ok := rdx.GetLastVal(data.VideoPublishTimeTextProperty, videoId); ok && ptts != "" {
				publishedDate = ptts
			}
		}

		if publishedDate != "" {
			properties[data.VideoPublishDateProperty] = publishedDate
		}
	}

	if len(properties) < 2 {
		if dts, ok := rdx.GetLastVal(data.VideoDownloadCompletedProperty, videoId); ok && dts != "" {
			properties[data.VideoDownloadCompletedProperty] = parseAndFormatDate(dts)
		}
	}

	return properties
}

func parseAndFormat(ts string) string {
	if pt, err := time.Parse(time.RFC3339, ts); err == nil {
		return pt.Local().Format(time.RFC1123)
	} else {
		return ts
	}
}

func parseAndFormatDate(ts string) string {
	if pt, err := time.Parse(time.RFC3339, ts); err == nil {
		return pt.Local().Format("Mon, 2 Jan 2006")
	} else {
		return ts
	}
}

func formatSeconds(ts int64) string {
	if ts == 0 {
		return "unknown"
	}

	t := time.Unix(ts, 0).UTC()

	layout := "4:05"
	if t.Hour() > 0 {
		layout = "15:04:05"
	}

	return t.Format(layout)
}

func channelTile(channelId string, rdx redux.Readable) strom.Element {

	tileContainer := strom.Create("a", "d-f", "fd-c", "br-s", "p-n").
		SetAttribute("href", path.Join("/channel", channelId)).
		SetStyle(map[string]string{
			"flow-shrink": "0",
			"padding":     "calc(1.5 * " + vars.Size(vars.SizeSmall) + ")",
			"row-gap":     vars.Size(vars.SizeXXSmall),
			"background":  vars.Color(vars.ColorHighlight),
			"width":       "max-content",
		})

	var title string
	if tp, ok := rdx.GetLastVal(data.ChannelTitleProperty, channelId); ok && tp != "" {
		title = tp
	}

	tileContainer.Append(strom.CreateText("span", title, "fw-b"))

	var newSubtitle string
	cnev := yeti.ChannelNotEndedVideos(channelId, math.MaxInt, rdx)
	if len(cnev) > 0 {
		switch len(cnev) {
		case 1:
			newSubtitle = "1 new video"
		default:
			newSubtitle = strconv.Itoa(len(cnev)) + " new videos"
		}
	} else {
		newSubtitle = "No new videos"
	}

	tileContainer.Append(strom.CreateText("span", newSubtitle).
		SetStyle(map[string]string{
			"font-size": vars.FontSize(vars.SizeXSmall),
			"color":     vars.Color(vars.ColorGray),
		}))

	return tileContainer
}

var extractedSearchVideosProperties = []string{
	data.VideoTitleProperty,
	data.VideoOwnerChannelNameProperty,
	data.VideoViewCountProperty,
	data.VideoPublishTimeTextProperty,
	data.VideoEndedDateProperty,
}

var extractedSearchPlaylistProperties = []string{
	data.PlaylistTitleProperty,
	data.PlaylistChannelProperty,
}

var extractedSearchChannelProperties = []string{
	data.ChannelTitleProperty,
	data.ChannelDescriptionProperty,
}

func extractSearchVideosMetadata(svrs []youtube_urls.VideoRenderer) map[string]map[string][]string {
	pkv := make(map[string]map[string][]string)

	for _, property := range extractedSearchVideosProperties {

		pkv[property] = make(map[string][]string)

		for _, svr := range svrs {

			id := svr.VideoId

			switch property {
			case data.VideoTitleProperty:
				pkv[property][id] = []string{svr.Title.String()}
			case data.VideoOwnerChannelNameProperty:
				pkv[property][id] = []string{svr.OwnerText.String()}
			case data.VideoViewCountProperty:
				pkv[property][id] = []string{svr.ViewCountText.SimpleText}
			case data.VideoPublishTimeTextProperty:
				pkv[property][id] = []string{svr.PublishedTimeText.SimpleText}
			}
		}

	}

	return pkv
}

func extractSearchPlaylistMetadata(sprs []youtube_urls.PlaylistRenderer) map[string]map[string][]string {
	pkv := make(map[string]map[string][]string)

	for _, property := range extractedSearchPlaylistProperties {

		pkv[property] = make(map[string][]string)

		for _, pvr := range sprs {

			id := pvr.PlaylistId

			switch property {
			case data.PlaylistTitleProperty:
				pkv[property][id] = []string{pvr.Title.SimpleText}
			}
		}

	}

	return pkv
}

func extractSearchChannelMetadata(scrs []youtube_urls.ChannelRenderer) map[string]map[string][]string {
	pkv := make(map[string]map[string][]string)

	for _, property := range extractedSearchChannelProperties {

		pkv[property] = make(map[string][]string)

		for _, cr := range scrs {

			id := cr.ChannelId

			switch property {
			case data.ChannelTitleProperty:
				pkv[property][id] = []string{cr.Title.SimpleText}
			case data.ChannelDescriptionProperty:
				pkv[property][id] = []string{cr.DescriptionSnippet.String()}
			}
		}

	}

	return pkv
}
