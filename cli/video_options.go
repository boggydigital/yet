package cli

import (
	"github.com/boggydigital/redux"
	"github.com/boggydigital/yet/data"
)

type VideoOptions struct {
	Favorite      bool
	DownloadQueue bool
	Progress      bool
	Ended         bool
	BgUtilBaseUrl string
	Reason        data.VideoEndedReason
	Verbose       bool
	Force         bool
}

func DefaultVideoOptions() *VideoOptions {
	return &VideoOptions{
		Favorite:      false,
		DownloadQueue: false,
		Progress:      false,
		Ended:         false,
		Reason:        data.DefaultEndedReason,
		Verbose:       false,
		Force:         false,
	}
}

func ApplyVideoDownloadOptions(opt *VideoOptions, videoId string, rdx redux.Readable) *VideoOptions {
	if f, ok := rdx.GetLastVal(data.VideoForcedDownloadProperty, videoId); ok && f == data.TrueValue {
		opt.Force = true
	}
	return opt
}
