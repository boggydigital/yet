package view_models

import (
	"fmt"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/yet/data"
	"strings"
)

const (
	showImagesLimit = 20
)

type PlaylistViewModel struct {
	PlaylistId      string
	PlaylistTitle   string
	PlaylistClass   string
	NewVideos       int
	Watching        bool
	AutoDownloading bool
	Videos          []*VideoViewModel
}

func GetPlaylistViewModel(playlistId string, rdx kvas.ReadableRedux) *PlaylistViewModel {

	nvc := 0

	if nv, ok := rdx.GetAllValues(data.PlaylistNewVideosProperty, playlistId); ok {
		nvc = len(nv)
	}

	watching := false
	if pwl, ok := rdx.GetFirstVal(data.PlaylistWatchlistProperty, playlistId); ok && pwl == data.TrueValue {
		watching = true
	}

	autoDownloading := false
	if pdq, ok := rdx.GetFirstVal(data.PlaylistDownloadQueueProperty, playlistId); ok && pdq == data.TrueValue {
		autoDownloading = true
	}

	pc := ""
	if autoDownloading {
		pc = "autoDownloading"
		if nvc == 0 {
			pc += " ended"
		}
	}

	plvm := &PlaylistViewModel{
		PlaylistId:      playlistId,
		PlaylistClass:   pc,
		NewVideos:       nvc,
		PlaylistTitle:   PlaylistTitle(playlistId, rdx),
		Watching:        watching,
		AutoDownloading: autoDownloading,
	}

	if videoIds, ok := rdx.GetAllValues(data.PlaylistVideosProperty, playlistId); ok && len(videoIds) > 0 {
		for i, videoId := range videoIds {
			var options []VideoOptions
			if i+1 < showImagesLimit {
				options = []VideoOptions{ShowPoster, ShowPublishedDate}
			} else {
				options = []VideoOptions{ShowPublishedDate}
			}
			plvm.Videos = append(plvm.Videos, GetVideoViewModel(videoId, rdx, options...))
		}
	}
	return plvm
}

func PlaylistTitle(playlistId string, rdx kvas.ReadableRedux) string {
	if plt, ok := rdx.GetFirstVal(data.PlaylistTitleProperty, playlistId); ok && plt != "" {

		if plc, ok := rdx.GetFirstVal(data.PlaylistChannelProperty, playlistId); ok && plc != "" && !strings.Contains(plt, plc) {
			if plt == "Videos" {
				return plc
			} else {
				return fmt.Sprintf("%s | %s", plc, plt)
			}
		}

		return plt
	}

	return playlistId
}
