package compton_elements

import (
	"embed"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/consts/size"
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/yet/data"
	"strconv"
	"time"
)

//go:embed "styles/*.css"
var videoLinkStyles embed.FS

type VideoDisplayOptions struct {
	Duration    bool
	PublishDate bool
	EndedDate   bool
	Downloaded  bool
}

func DefaultVideoDisplayOptions() *VideoDisplayOptions {
	return &VideoDisplayOptions{
		Duration:    true,
		PublishDate: true,
		EndedDate:   false,
		Downloaded:  false,
	}
}

func VideoLink(r compton.Registrar, videoId string, rdx kevlar.ReadableRedux, opt *VideoDisplayOptions) compton.Element {

	if opt == nil {
		opt = DefaultVideoDisplayOptions()
	}

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

	if opt.Duration {
		if durs, sure := rdx.GetLastVal(data.VideoDurationProperty, videoId); sure && durs != "" {
			if duri, err := strconv.ParseInt(durs, 10, 64); err == nil {
				durationDiv := compton.DivText(formatSeconds(duri))
				durationDiv.AddClass("duration")
				stack.Append(durationDiv)
			}
		}
	}

	if title, ok := rdx.GetLastVal(data.VideoTitleProperty, videoId); ok {
		stack.Append(compton.H2Text(title))
	}

	if opt.PublishDate {
		var publishedDate string
		if pds, ok := rdx.GetLastVal(data.VideoPublishDateProperty, videoId); ok && pds != "" {
			publishedDate = parseAndFormat(pds)
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

	if opt.Downloaded {
		downFrow := compton.Frow(r).FontSize(size.Small)
		if dts, ok := rdx.GetLastVal(data.VideoDownloadCompletedProperty, videoId); ok && dts != "" {
			downFrow.PropVal("Downloaded", parseAndFormat(dts))
			stack.Append(downFrow)
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
