package view_models

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet_urls/youtube_urls"
	"strings"
)

type VideoManagementViewModel struct {
	VideoId         string
	VideoTitle      string
	CanViewAtOrigin bool
	Progress        bool
	CurrentTime     string
	Ended           bool
	Skipped         bool
	Watchlist       bool
	DownloadQueue   bool
	ForcedDownload  bool
	SingleFormat    bool
}

func GetVideoManagementModel(videoId string, rdx kvas.ReadableRedux) *VideoManagementViewModel {
	videoTitle := ""
	if vt, ok := rdx.GetLastVal(data.VideoTitleProperty, videoId); ok && vt != "" {
		videoTitle = vt
	}

	currentTime := ""
	if ct, ok := rdx.GetLastVal(data.VideoProgressProperty, videoId); ok && ct != "" {
		currentTime = ct
	}

	return &VideoManagementViewModel{
		VideoId:         videoId,
		VideoTitle:      videoTitle,
		CanViewAtOrigin: !strings.Contains(videoId, youtube_urls.DefaultVideoExt),
		CurrentTime:     currentTime,
		Progress:        rdx.HasKey(data.VideoProgressProperty, videoId),
		Ended:           rdx.HasKey(data.VideoEndedDateProperty, videoId),
		//TODO: update to reason
		Skipped: rdx.HasKey(data.VideoEndedReasonProperty, videoId),
		//Watchlist:      rdx.HasKey(data.VideosWatchlistProperty, videoId),
		DownloadQueue:  rdx.HasKey(data.VideoDownloadQueuedProperty, videoId),
		ForcedDownload: rdx.HasKey(data.VideoForcedDownloadProperty, videoId),
		SingleFormat:   rdx.HasKey(data.VideoPreferSingleFormatProperty, videoId),
	}
}
