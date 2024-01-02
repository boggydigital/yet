package rest

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/yet/data"
	"slices"
	"strconv"
	"time"
)

type VideoOptions int

const (
	ShowPoster VideoOptions = iota
	ShowPublishedDate
	ShowEndedDate
	ShowRemainingDuration
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

	optShowPublishedDate := slices.Contains(options, ShowPublishedDate)

	if optShowPublishedDate {
		if pts, ok := rdx.GetFirstVal(data.VideoPublishDateProperty, videoId); ok && pts != "" {
			publishedDate = parseAndFormat(pts)
		}
		if dts, ok := rdx.GetFirstVal(data.VideoDownloadedDateProperty, videoId); ok && dts != "" {
			downloadedDate = parseAndFormat(dts)
		}
	}

	optShowEndedDate := slices.Contains(options, ShowEndedDate)

	ended := false
	endedDate := ""
	if ets, ok := rdx.GetFirstVal(data.VideoEndedProperty, videoId); ok && ets != "" {
		ended = true
		if optShowEndedDate {
			endedDate = parseAndFormat(ets)
		}
	}

	var rem, dur int64

	optShowRemainingDuration := slices.Contains(options, ShowRemainingDuration)

	if optShowRemainingDuration {
		var ct int64
		if cts, ok := rdx.GetFirstVal(data.VideoProgressProperty, videoId); ok && cts != "" {
			if cti, err := strconv.ParseInt(cts, 10, 64); err == nil {
				ct = cti
			}
		}
		if durs, sure := rdx.GetFirstVal(data.VideoDurationProperty, videoId); sure && durs != "" {
			if duri, err := strconv.ParseInt(durs, 10, 64); err == nil {
				dur = duri
			}
		}
		rem = dur - ct
	}

	class := ""
	if ended {
		videoTitle = "☑️ " + videoTitle
		class += "ended"
	}

	optShowPoster := slices.Contains(options, ShowPoster)

	return &VideoViewModel{
		VideoId:               videoId,
		VideoUrl:              videoUrl,
		VideoTitle:            videoTitle,
		Class:                 class,
		ShowPoster:            optShowPoster,
		ShowPublishedDate:     optShowPublishedDate,
		PublishedDate:         publishedDate,
		DownloadedDate:        downloadedDate,
		ShowEndedDate:         optShowEndedDate,
		EndedDate:             endedDate,
		ShowRemainingDuration: optShowRemainingDuration && dur > 0,
		RemainingTime:         formatSeconds(rem),
		Duration:              formatSeconds(dur),
	}

}

func parseAndFormat(ts string) string {
	if pt, err := time.Parse(time.RFC3339, ts); err == nil {
		return pt.Local().Format(time.RFC1123)
	} else {
		return ts
	}
}

func formatSeconds(ts int64) string {
	dur := time.Duration(float64(ts) * float64(time.Second))
	return dur.String()
}
