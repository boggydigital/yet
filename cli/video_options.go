package cli

import (
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/yet/data"
)

type VideoOptions struct {
	Favorite      bool
	DownloadQueue bool
	Progress      bool
	Ended         bool
	Reason        data.VideoEndedReason
	Source        string
	Force         bool
}

func DefaultVideoOptions() *VideoOptions {
	return &VideoOptions{
		Favorite:      false,
		DownloadQueue: false,
		Progress:      false,
		Ended:         false,
		Reason:        data.DefaultEndedReason,
		Source:        "",
		Force:         false,
	}
}

func ApplyVideoDownloadOptions(opt *VideoOptions, videoId string, rdx kevlar.ReadableRedux) *VideoOptions {
	if f, ok := rdx.GetLastVal(data.VideoForcedDownloadProperty, videoId); ok && f == data.TrueValue {
		opt.Force = true
	}
	if src, ok := rdx.GetLastVal(data.VideoSourceProperty, videoId); ok && src != "" {
		opt.Source = src
	}
	return opt
}
