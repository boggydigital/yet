package rest

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/yet/data"
	"net/http"
	"slices"
	"strings"
	"time"
)

type VideoOptions int

const (
	ShowPoster VideoOptions = iota
	ShowPublishedDate
	ShowEndedDate
)

func writeVideo(videoId string, rdx kvas.ReadableRedux, sb *strings.Builder, options ...VideoOptions) {

	videoTitle := videoId
	if title, ok := rdx.GetFirstVal(data.VideoTitleProperty, videoId); ok && title != "" {
		videoTitle = title
	}

	videoUrl := "/watch?"
	if videoId != "" {
		videoUrl += "v=" + videoId
	}

	posterContent := ""
	if slices.Contains(options, ShowPoster) {
		posterContent = "<img src='/poster?v=" + videoId + "&q=mqdefault' loading='lazy'/>"
	}

	publishedContent := ""
	if slices.Contains(options, ShowPublishedDate) {
		if pts, ok := rdx.GetFirstVal(data.VideoPublishDateProperty, videoId); ok && pts != "" {
			publishedContent = "<span class='subtitle'><b>Published</b>: " + parseAndFormat(time.RFC3339, pts) + "</span>"
		}
	}

	ended := false
	endedContent := ""
	if ets, ok := rdx.GetFirstVal(data.VideoEndedProperty, videoId); ok && ets != "" {
		ended = true
		if slices.Contains(options, ShowEndedDate) {
			endedContent = "<span class='subtitle'><b>Ended</b>: " + parseAndFormat(http.TimeFormat, ets) + "</span>"
		}
	}

	class := ""
	if ended {
		videoTitle = "☑️ " + videoTitle
		class += "ended"
	}

	sb.WriteString("<a class='" + class + "' href='" + videoUrl + "'>" +
		posterContent +
		"<span class='title'>" + videoTitle + "</span>" +
		endedContent +
		publishedContent +
		"</a>")
}

func parseAndFormat(layout, ts string) string {
	if pt, err := time.Parse(layout, ts); err == nil {
		return pt.Local().Format(time.RFC1123)
	} else {
		return ts
	}
}
