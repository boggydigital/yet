package compton_elements

import (
	"embed"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/color"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/consts/font_weight"
	"github.com/boggydigital/compton/consts/size"
	"github.com/boggydigital/redux"
	"github.com/boggydigital/yet/data"
	"strconv"
	"time"
)

//go:embed "styles/*.css"
var videoLinkStyles embed.FS

var propertyTitles = map[string]string{
	data.VideoOwnerChannelNameProperty:  "Channel",
	data.VideoEndedDateProperty:         "Ended",
	data.VideoPublishDateProperty:       "Published",
	data.VideoDownloadCompletedProperty: "Downloaded",
	data.VideoEndedReasonProperty:       "How",
}

var propertiesOrder = []string{
	data.VideoOwnerChannelNameProperty,
	data.VideoEndedDateProperty,
	data.VideoPublishDateProperty,
	data.VideoDownloadCompletedProperty,
	data.VideoEndedReasonProperty,
}

func VideoLink(r compton.Registrar, videoId string, rdx redux.Readable) compton.Element {

	r.RegisterStyles(videoLinkStyles, "styles/video-link.css")

	link := compton.A("/watch?v=" + videoId)
	link.AddClass("video-link")

	if rdx.HasKey(data.VideoEndedDateProperty, videoId) {
		link.AddClass("ended")
	}

	stack := compton.FlexItems(r, direction.Column).RowGap(size.Small)
	link.Append(stack)

	var dehydratedPoster string
	var repColor string

	if dhp, ok := rdx.GetLastVal(data.VideoDehydratedPosterProperty, videoId); ok {
		dehydratedPoster = dhp
	}

	if rc, ok := rdx.GetLastVal(data.VideoDehydratedRepColorProperty, videoId); ok {
		repColor = rc
	}

	issaImage := compton.IssaImageDehydrated(r, repColor, dehydratedPoster, "/poster?v="+videoId+"&q=hqdefault")
	issaImage.AspectRatio(float64(16) / float64(9))
	issaImage.AddClass("poster")

	stack.Append(issaImage)

	if durs, sure := rdx.GetLastVal(data.VideoDurationProperty, videoId); sure && durs != "" {
		if duri, err := strconv.ParseInt(durs, 10, 64); err == nil {

			var remaining int64

			if cts, ok := rdx.GetLastVal(data.VideoProgressProperty, videoId); ok && cts != "" {
				if cti, err := strconv.ParseInt(cts, 10, 64); err == nil {
					remaining = duri - cti
				}
			}

			durationItems := compton.FlexItems(r, direction.Row).FontSize(size.Small).ColumnGap(size.Small)
			durationItems.AddClass("duration")

			durSpan := compton.Fspan(r, formatSeconds(duri))

			if remaining > 0 {
				remSpan := compton.Fspan(r, formatSeconds(remaining)).FontWeight(font_weight.Bolder)
				durationItems.Append(remSpan)
				durSpan.ForegroundColor(color.Gray)
			} else {
				durSpan.FontWeight(font_weight.Bolder)
			}

			durationItems.Append(durSpan)

			stack.Append(durationItems)
		}
	}

	if title, ok := rdx.GetLastVal(data.VideoTitleProperty, videoId); ok {
		stack.Append(compton.H2Text(title))
	}

	vsp := videoSummaryProperties(videoId, rdx)

	for _, p := range propertiesOrder {
		v := vsp[p]
		if v == "" {
			continue
		}
		fr := compton.Frow(r).FontSize(size.Small)
		fr.PropVal(propertyTitles[p], v)

		if p == data.VideoEndedDateProperty {
			endedReason := data.DefaultEndedReason
			if ers, ok := rdx.GetLastVal(data.VideoEndedReasonProperty, videoId); ok {
				endedReason = data.ParseVideoEndedReason(ers)
			}
			fr.PropVal(propertyTitles[data.VideoEndedReasonProperty], string(endedReason))
		}

		stack.Append(fr)
	}

	return link
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
