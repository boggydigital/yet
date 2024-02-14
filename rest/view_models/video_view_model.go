package view_models

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
	ShowDuration
	ShowOwnerChannel
	ShowViewCount
)

type VideoViewModel struct {
	VideoId               string
	VideoUrl              string
	VideoTitle            string
	Class                 string
	ShowPoster            bool
	ShowPublishedDate     bool
	PublishedDate         string
	DownloadedDate        string
	ShowEndedDate         bool
	EndedDate             string
	RemainingTime         string
	Duration              string
	ShowDuration          bool
	ShowRemainingDuration bool
	CurrentTimeSeconds    string
	DurationSeconds       string
	ShowOwnerChannel      bool
	OwnerChannel          string
	ShowViewCount         bool
	ViewCount             string
}

func GetVideoViewModel(videoId string, rdx kvas.ReadableRedux, options ...VideoOptions) *VideoViewModel {

	videoTitle := videoId
	if title, ok := rdx.GetLastVal(data.VideoTitleProperty, videoId); ok && title != "" {
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
		if pds, ok := rdx.GetLastVal(data.VideoPublishDateProperty, videoId); ok && pds != "" {
			publishedDate = parseAndFormat(pds)
		} else {
			if ptts, ok := rdx.GetLastVal(data.VideoPublishTimeTextProperty, videoId); ok && ptts != "" {
				publishedDate = ptts
			}
		}
		if dts, ok := rdx.GetLastVal(data.VideoDownloadedDateProperty, videoId); ok && dts != "" {
			downloadedDate = parseAndFormat(dts)
		}
	}

	ownerChannel := ""

	optShowOwnerChannel := slices.Contains(options, ShowOwnerChannel)
	if optShowOwnerChannel {
		if och, ok := rdx.GetLastVal(data.VideoOwnerChannelNameProperty, videoId); ok && och != "" {
			ownerChannel = och
		}
	}

	optShowEndedDate := slices.Contains(options, ShowEndedDate)

	ended := false
	endedDate := ""
	if ets, ok := rdx.GetLastVal(data.VideoEndedProperty, videoId); ok && ets != "" {
		ended = true
		if optShowEndedDate {
			endedDate = parseAndFormat(ets)
		}
	}

	skipped := false
	if s, ok := rdx.GetLastVal(data.VideoSkippedProperty, videoId); ok && s != "" {
		skipped = true
	}

	var rem, dur int64

	optShowDuration := slices.Contains(options, ShowDuration)
	optShowRemainingDuration := slices.Contains(options, ShowRemainingDuration)

	if optShowDuration || optShowRemainingDuration {
		if durs, sure := rdx.GetLastVal(data.VideoDurationProperty, videoId); sure && durs != "" {
			if duri, err := strconv.ParseInt(durs, 10, 64); err == nil {
				dur = duri
			}
		}
	}

	if optShowRemainingDuration {
		var ct int64
		if cts, ok := rdx.GetLastVal(data.VideoProgressProperty, videoId); ok && cts != "" {
			if cti, err := strconv.ParseInt(cts, 10, 64); err == nil {
				ct = cti
			}
		}
		rem = dur - ct
	}

	viewCount := ""
	optShowViewCount := slices.Contains(options, ShowViewCount)

	if optShowViewCount {
		if vc, ok := rdx.GetLastVal(data.VideoViewCountProperty, videoId); ok && vc != "" {
			viewCount = vc
		}
	}

	class := ""
	if ended {
		titlePrefix := "☑️ "
		class += "ended "
		if skipped {
			titlePrefix = "⏭️ "
			class += "skipped "
		}
		videoTitle = titlePrefix + videoTitle
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
		ShowDuration:          optShowDuration,
		Duration:              formatSeconds(dur),
		ShowRemainingDuration: optShowRemainingDuration && dur > 0,
		RemainingTime:         formatSeconds(rem),
		CurrentTimeSeconds:    strconv.FormatInt(dur-rem, 10),
		DurationSeconds:       strconv.FormatInt(dur, 10),
		ShowOwnerChannel:      optShowOwnerChannel,
		OwnerChannel:          ownerChannel,
		ShowViewCount:         optShowViewCount,
		ViewCount:             viewCount,
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
	if ts == 0 {
		return "unknown"
	}
	dur := time.Duration(float64(ts) * float64(time.Second))
	return dur.String()
}
