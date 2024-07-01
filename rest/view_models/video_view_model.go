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
	ShowProgress
	ShowDuration
	ShowOwnerChannel
	ShowViewCount
)

type VideoViewModel struct {
	VideoId            string
	VideoUrl           string
	VideoTitle         string
	Favorite           bool
	ShowPoster         bool
	ShowPublishedDate  bool
	PublishedDate      string
	DownloadedDate     string
	ShowEndedTime      bool
	EndedTime          string
	EndedReason        data.VideoEndedReason
	ShowDuration       bool
	Duration           string
	ShowProgress       bool
	CurrentTimeSeconds string
	DurationSeconds    string
	ShowOwnerChannel   bool
	OwnerChannel       string
	ShowViewCount      bool
	ViewCount          string
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
		if dts, ok := rdx.GetLastVal(data.VideoDownloadCompletedProperty, videoId); ok && dts != "" {
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

	endedDate := ""
	if ets, ok := rdx.GetLastVal(data.VideoEndedDateProperty, videoId); ok && ets != "" {
		endedDate = parseAndFormat(ets)
	}

	var rem, dur int64

	optShowDuration := slices.Contains(options, ShowDuration)
	optShowProgress := slices.Contains(options, ShowProgress)

	if optShowDuration || optShowProgress {
		if durs, sure := rdx.GetLastVal(data.VideoDurationProperty, videoId); sure && durs != "" {
			if duri, err := strconv.ParseInt(durs, 10, 64); err == nil {
				dur = duri
			}
		}
	}

	if optShowProgress {
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

	endedReason := data.DefaultEndedReason
	if er, ok := rdx.GetLastVal(data.VideoEndedReasonProperty, videoId); ok {
		endedReason = data.ParseVideoEndedReason(er)
	}

	favorite := false
	if rdx.HasKey(data.VideoFavoriteProperty, videoId) {
		favorite = true
	}

	optShowPoster := slices.Contains(options, ShowPoster)

	return &VideoViewModel{
		VideoId:            videoId,
		VideoUrl:           videoUrl,
		VideoTitle:         videoTitle,
		Favorite:           favorite,
		ShowPoster:         optShowPoster,
		ShowPublishedDate:  optShowPublishedDate,
		PublishedDate:      publishedDate,
		DownloadedDate:     downloadedDate,
		ShowEndedTime:      optShowEndedDate,
		EndedTime:          endedDate,
		EndedReason:        endedReason,
		ShowDuration:       optShowDuration && dur > 0,
		Duration:           formatSeconds(dur),
		ShowProgress:       optShowProgress && rem > 0 && dur > 0,
		CurrentTimeSeconds: strconv.FormatInt(dur-rem, 10),
		DurationSeconds:    strconv.FormatInt(dur, 10),
		ShowOwnerChannel:   optShowOwnerChannel,
		OwnerChannel:       ownerChannel,
		ShowViewCount:      optShowViewCount,
		ViewCount:          viewCount,
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

	t := time.Unix(ts, 0).UTC()

	layout := "4:05"
	if t.Hour() > 0 {
		layout = "15:04:05"
	}

	return t.Format(layout)
}
