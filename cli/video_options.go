package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/yet/data"
)

type VideoOptions struct {
	Favorite           bool
	DownloadQueue      bool
	Progress           bool
	Ended              bool
	Reason             data.VideoEndedReason
	Source             string
	PreferSingleFormat bool
	Force              bool
}

func DefaultVideoOptions() *VideoOptions {
	return &VideoOptions{
		Favorite:           false,
		DownloadQueue:      false,
		Progress:           false,
		Ended:              false,
		Reason:             data.DefaultEndedReason,
		Source:             "",
		PreferSingleFormat: true,
		Force:              false,
	}
}

func ApplyVideoDownloadOptions(opt *VideoOptions, videoId string, rdx kvas.ReadableRedux) *VideoOptions {
	if f, ok := rdx.GetLastVal(data.VideoForcedDownloadProperty, videoId); ok && f == data.TrueValue {
		opt.Force = true
	}
	if psf, ok := rdx.GetLastVal(data.VideoPreferSingleFormatProperty, videoId); ok && psf == data.TrueValue {
		opt.PreferSingleFormat = true
	}
	if src, ok := rdx.GetLastVal(data.VideoSourceProperty, videoId); ok && src != "" {
		opt.Source = src
	}
	return opt
}
