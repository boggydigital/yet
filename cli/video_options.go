package cli

import "github.com/boggydigital/yet/data"

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
