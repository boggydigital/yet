package cli

import "github.com/boggydigital/yet/data"

type PlaylistOptions struct {
	AutoRefresh        bool
	AutoDownload       bool
	DownloadPolicy     data.PlaylistDownloadPolicy
	PreferSingleFormat bool
	Expand             bool
	Force              bool
}

func DefaultPlaylistOptions() *PlaylistOptions {
	return &PlaylistOptions{
		AutoRefresh:        false,
		AutoDownload:       false,
		DownloadPolicy:     data.Recent,
		PreferSingleFormat: true,
		Expand:             false,
		Force:              false,
	}
}
