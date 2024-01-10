package view_models

import (
	"fmt"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/yet/data"
	"strings"
)

const (
	showImagesLimit = 12
)

type PlaylistViewModel struct {
	PlaylistId    string
	PlaylistTitle string
	PlaylistClass string
	NewVideos     int
	Watching      bool
	Downloading   bool
	Videos        []*VideoViewModel
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

	downloading := false
	if pdq, ok := rdx.GetFirstVal(data.PlaylistDownloadQueueProperty, playlistId); ok && pdq == data.TrueValue {
		downloading = true
	}

	pc := ""
	if downloading {
		pc = "downloading"
		if nvc == 0 {
			pc += " ended"
		}
	}

	plvm := &PlaylistViewModel{
		PlaylistId:    playlistId,
		PlaylistClass: pc,
		NewVideos:     nvc,
		PlaylistTitle: PlaylistTitle(playlistId, rdx),
		Watching:      watching,
		Downloading:   downloading,
	}

	if videoIds, ok := rdx.GetAllValues(data.PlaylistVideosProperty, playlistId); ok && len(videoIds) > 0 {
		for i, videoId := range videoIds {
			var options []VideoOptions
			if i+1 < showImagesLimit {
				options = []VideoOptions{ShowPoster, ShowViewCount, ShowPublishedDate}
			} else {
				options = []VideoOptions{ShowViewCount, ShowPublishedDate}
			}
			plvm.Videos = append(plvm.Videos, GetVideoViewModel(videoId, rdx, options...))
		}
	}
	return plvm
}

func PlaylistTitle(playlistId string, rdx kvas.ReadableRedux) string {
	if plt, ok := rdx.GetFirstVal(data.PlaylistTitleProperty, playlistId); ok && plt != "" {

		if plc, ok := rdx.GetFirstVal(data.PlaylistChannelProperty, playlistId); ok && plc != "" && !strings.Contains(plt, plc) {
			//if plt == "Videos" {
			//	return plc
			//} else {
			return fmt.Sprintf("%s · %s", plt, plc)
			//}
		}

		return plt
	}

	return playlistId
}
