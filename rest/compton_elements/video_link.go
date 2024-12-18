package compton_elements

import (
	"embed"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/color"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/consts/size"
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/yet/data"
	"slices"
	"strconv"
	"time"
)

//go:embed "styles/*.css"
var videoLinkStyles embed.FS

type VideoDisplayOptions int

const (
	ShowPublishedDate VideoDisplayOptions = iota
	ShowDownloadedDate
	ShowEndedDate
	//ShowDuration
	ShowOwnerChannel
)

func VideoLink(r compton.Registrar, videoId string, rdx kevlar.ReadableRedux, options ...VideoDisplayOptions) compton.Element {

	r.RegisterStyles(videoLinkStyles, "styles/video-link.css")

	link := compton.A("/watch?v=" + videoId)
	link.AddClass("video-link")

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
				remSpan := compton.SpanText(formatSeconds(remaining))
				durationItems.Append(remSpan)
				durSpan.ForegroundColor(color.Gray)
			}

			durationItems.Append(durSpan)

			stack.Append(durationItems)
		}
	}

	if title, ok := rdx.GetLastVal(data.VideoTitleProperty, videoId); ok {
		stack.Append(compton.H2Text(title))
	}

	if slices.Contains(options, ShowOwnerChannel) {
		if och, ok := rdx.GetLastVal(data.VideoOwnerChannelNameProperty, videoId); ok && och != "" {
			channelFrow := compton.Frow(r).FontSize(size.Small)
			channelFrow.PropVal("Channel", och)
			stack.Append(channelFrow)
		}
	}

	if slices.Contains(options, ShowPublishedDate) {
		var publishedDate string
		if pds, ok := rdx.GetLastVal(data.VideoPublishDateProperty, videoId); ok && pds != "" {
			publishedDate = parseAndFormatDate(pds)
		} else {
			if ptts, ok := rdx.GetLastVal(data.VideoPublishTimeTextProperty, videoId); ok && ptts != "" {
				publishedDate = ptts
			}
		}

		if publishedDate != "" {
			pubFrow := compton.Frow(r).FontSize(size.Small)
			pubFrow.PropVal("Published", publishedDate)
			stack.Append(pubFrow)
		}
	}

	if slices.Contains(options, ShowDownloadedDate) {
		if dts, ok := rdx.GetLastVal(data.VideoDownloadCompletedProperty, videoId); ok && dts != "" {
			downFrow := compton.Frow(r).FontSize(size.Small)
			downFrow.PropVal("Downloaded", parseAndFormatDate(dts))
			stack.Append(downFrow)
		}
	}

	if ets, ok := rdx.GetLastVal(data.VideoEndedDateProperty, videoId); ok && ets != "" {
		link.AddClass("ended")
		if slices.Contains(options, ShowEndedDate) {
			endedFrow := compton.Frow(r).FontSize(size.Small)
			endedFrow.PropVal("Ended", parseAndFormatDate(ets))

			endedReason := data.DefaultEndedReason
			if er, ok := rdx.GetLastVal(data.VideoEndedReasonProperty, videoId); ok {
				endedReason = data.ParseVideoEndedReason(er)
			}
			endedFrow.PropVal("How", endedReason.String())

			stack.Append(endedFrow)
		}
	}

	return link
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
