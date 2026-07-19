package rest

import (
	"iter"
	"net/http"
	"strings"

	"github.com/boggydigital/redux"
	"github.com/boggydigital/strom"
	"github.com/boggydigital/strom/vars/atoms"
	"github.com/boggydigital/strom/vars/sizes"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/rest/view_models"
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

	root, body := strom.RootBody("Search", atoms.FlexCol(sizes.Normal)...)

	topRow := strom.Create("ul", atoms.FlexRow(sizes.Small)...).AddAtom(atoms.AlignItemsCenter)
	body.Append(topRow)

	topRow.Append(navButton("Home", "/"))
	topRow.Append(strom.CreateText("h2", "Results for '"+searchQuery+"'"))

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

		channels := strom.Create("ul", atoms.FlexRowWrap(sizes.Normal)...)
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

		videos := strom.Create("ul", atoms.FlexRowWrap(sizes.Normal)...)
		body.Append(videos)

		srv := new(searchResultsVideos{searchInitialData: sid, rdx: rdx})

		videos.Append(strom.OnDemand(srv.getVideos))
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
