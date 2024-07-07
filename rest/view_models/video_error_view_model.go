package view_models

import (
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/yet/data"
)

type VideoErrorViewModel struct {
	VideoId    string
	VideoTitle string
	Error      string
}

func GetVideoErrorViewModel(videoId, error string, rdx kevlar.ReadableRedux) *VideoErrorViewModel {
	videoTitle := ""
	if vt, ok := rdx.GetLastVal(data.VideoTitleProperty, videoId); ok && vt != "" {
		videoTitle = vt
	}

	if videoTitle == "" {
		videoTitle = videoId
	}

	return &VideoErrorViewModel{
		VideoId:    videoId,
		VideoTitle: videoTitle,
		Error:      error,
	}
}
