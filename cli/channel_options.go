package cli

import "github.com/boggydigital/yet/data"

type ChannelOptions struct {
	Playlists          bool
	AutoRefresh        bool
	AutoDownload       bool
	DownloadPolicy     data.DownloadPolicy
	PreferSingleFormat bool
	Expand             bool
	Force              bool
}

func DefaultChannelOptions() *ChannelOptions {
	return &ChannelOptions{
		AutoRefresh:        false,
		AutoDownload:       false,
		DownloadPolicy:     data.Recent,
		PreferSingleFormat: true,
		Expand:             false,
		Force:              false,
	}
}
