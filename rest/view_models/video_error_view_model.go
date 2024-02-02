package view_models

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/yet/data"
)

type VideoErrorViewModel struct {
	VideoId    string
	VideoTitle string
	Error      string
}

func GetVideoErrorViewModel(videoId, error string, rdx kvas.ReadableRedux) *VideoErrorViewModel {
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
