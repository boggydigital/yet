package view_models

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yt_urls"
	"strings"
)

type VideoManagementViewModel struct {
	VideoId         string
	VideoTitle      string
	CanViewAtOrigin bool
	Progress        bool
	CurrentTime     string
	Ended           bool
	Watchlist       bool
	DownloadQueue   bool
}

func GetVideoManagementModel(videoId string, rdx kvas.ReadableRedux) *VideoManagementViewModel {
	videoTitle := ""
	if vt, ok := rdx.GetFirstVal(data.VideoTitleProperty, videoId); ok && vt != "" {
		videoTitle = vt
	}

	currentTime := ""
	if ct, ok := rdx.GetFirstVal(data.VideoProgressProperty, videoId); ok && ct != "" {
		currentTime = ct
	}

	return &VideoManagementViewModel{
		VideoId:         videoId,
		VideoTitle:      videoTitle,
		CanViewAtOrigin: !strings.Contains(videoId, yt_urls.DefaultVideoExt),
		CurrentTime:     currentTime,
		Progress:        rdx.HasKey(data.VideoProgressProperty, videoId),
		Ended:           rdx.HasKey(data.VideoEndedProperty, videoId),
		Watchlist:       rdx.HasKey(data.VideosWatchlistProperty, videoId),
		DownloadQueue:   rdx.HasKey(data.VideosDownloadQueueProperty, videoId),
	}
}
