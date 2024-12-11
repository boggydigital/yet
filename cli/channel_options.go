package cli

import "github.com/boggydigital/yet/data"

type ChannelOptions struct {
	Playlists      bool
	AutoRefresh    bool
	AutoDownload   bool
	DownloadPolicy data.DownloadPolicy
	Expand         bool
	Force          bool
}

func DefaultChannelOptions() *ChannelOptions {
	return &ChannelOptions{
		AutoRefresh:    false,
		AutoDownload:   false,
		DownloadPolicy: data.Recent,
		Expand:         false,
		Force:          false,
	}
}
