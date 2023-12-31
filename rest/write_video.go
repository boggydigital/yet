package rest

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/yet/data"
	"slices"
	"time"
)

type VideoOptions int

const (
	ShowPoster VideoOptions = iota
	ShowPublishedDate
	ShowEndedDate
)

func videoViewModel(videoId string, rdx kvas.ReadableRedux, options ...VideoOptions) *VideoViewModel {

	videoTitle := videoId
	if title, ok := rdx.GetFirstVal(data.VideoTitleProperty, videoId); ok && title != "" {
		videoTitle = title
	}

	videoUrl := "/watch?"
	if videoId != "" {
		videoUrl += "v=" + videoId
	}

	publishedDate := ""
	downloadedDate := ""

	if slices.Contains(options, ShowPublishedDate) {
		if pts, ok := rdx.GetFirstVal(data.VideoPublishDateProperty, videoId); ok && pts != "" {
			publishedDate = parseAndFormat(pts)
		}
		if dts, ok := rdx.GetFirstVal(data.VideoDownloadedDateProperty, videoId); ok && dts != "" {
			downloadedDate = parseAndFormat(dts)
		}
	}

	ended := false
	endedDate := ""
	if ets, ok := rdx.GetFirstVal(data.VideoEndedProperty, videoId); ok && ets != "" {
		ended = true
		if slices.Contains(options, ShowEndedDate) {
			endedDate = parseAndFormat(ets)
		}
	}

	class := ""
	if ended {
		videoTitle = "☑️ " + videoTitle
		class += "ended"
	}

	return &VideoViewModel{
		VideoId:           videoId,
		VideoUrl:          videoUrl,
		VideoTitle:        videoTitle,
		Class:             class,
		ShowPoster:        slices.Contains(options, ShowPoster),
		ShowPublishedDate: slices.Contains(options, ShowPublishedDate),
		PublishedDate:     publishedDate,
		DownloadedDate:    downloadedDate,
		ShowEndedDate:     slices.Contains(options, ShowEndedDate),
		EndedDate:         endedDate,
	}

}

func parseAndFormat(ts string) string {
	if pt, err := time.Parse(time.RFC3339, ts); err == nil {
		return pt.Local().Format(time.RFC1123)
	} else {
		return ts
	}
}
