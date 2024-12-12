package view_models

import (
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/yet/data"
)

type VideoManagementViewModel struct {
	VideoId         string
	VideoTitle      string
	CurrentTime     string
	Favorite        bool
	Progress        bool
	Ended           bool
	EndedReason     data.VideoEndedReason
	AllEndedReasons []data.VideoEndedReason
	DownloadQueued  bool
	ForcedDownload  bool
}

func GetVideoManagementModel(videoId string, rdx kevlar.ReadableRedux) *VideoManagementViewModel {
	videoTitle := ""
	if vt, ok := rdx.GetLastVal(data.VideoTitleProperty, videoId); ok && vt != "" {
		videoTitle = vt
	}

	currentTime := ""
	if ct, ok := rdx.GetLastVal(data.VideoProgressProperty, videoId); ok && ct != "" {
		currentTime = ct
	}

	endedReason := data.DefaultEndedReason
	if er, ok := rdx.GetLastVal(data.VideoEndedReasonProperty, videoId); ok && er != "" {
		endedReason = data.ParseVideoEndedReason(er)
	}

	downloadQueued := false
	if dqs, ok := rdx.GetLastVal(data.VideoDownloadQueuedProperty, videoId); ok {
		downloadQueued = true
		if dcs, sure := rdx.GetLastVal(data.VideoDownloadCompletedProperty, videoId); sure {
			if dqs < dcs {
				downloadQueued = false
			}
		}
	}

	return &VideoManagementViewModel{
		VideoId:         videoId,
		VideoTitle:      videoTitle,
		CurrentTime:     currentTime,
		Favorite:        rdx.HasKey(data.VideoFavoriteProperty, videoId),
		Progress:        rdx.HasKey(data.VideoProgressProperty, videoId),
		Ended:           rdx.HasKey(data.VideoEndedDateProperty, videoId),
		EndedReason:     endedReason,
		AllEndedReasons: data.AllVideoEndedReasons(),
		DownloadQueued:  downloadQueued,
		ForcedDownload:  rdx.HasKey(data.VideoForcedDownloadProperty, videoId),
	}
}
