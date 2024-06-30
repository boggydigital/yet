package view_models

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/yet/data"
)

type PlaylistViewModel struct {
	PlaylistId           string
	PlaylistTitle        string
	PlaylistChannelTitle string
	PlaylistClass        string
	NewVideos            int
	Watching             bool
	Downloading          bool
	SingleFormat         bool
	Videos               []*VideoViewModel
}

func GetPlaylistViewModel(playlistId string, rdx kvas.ReadableRedux) *PlaylistViewModel {

	if playlistId == "" {
		return nil
	}

	nvc := 0

	//if nv, ok := rdx.GetAllValues(data.PlaylistNewVideosProperty, playlistId); ok {
	//	nvc = len(nv)
	//}

	watching := false
	if pwl, ok := rdx.GetLastVal(data.PlaylistAutoRefreshProperty, playlistId); ok && pwl == data.TrueValue {
		watching = true
	}

	downloading := false
	if pdq, ok := rdx.GetLastVal(data.PlaylistAutoDownloadProperty, playlistId); ok && pdq == data.TrueValue {
		downloading = true
	}

	singleFormat := false
	if psf, ok := rdx.GetLastVal(data.PlaylistPreferSingleFormatProperty, playlistId); ok && psf == data.TrueValue {
		singleFormat = true
	}

	pc := ""
	if downloading {
		pc += " downloading"
	}
	if nvc == 0 {
		pc += " ended"
	}

	playlistTitle := ""
	if plt, ok := rdx.GetLastVal(data.PlaylistTitleProperty, playlistId); ok && plt != "" {
		playlistTitle = plt
	}

	playlistChannelTitle := ""
	if plc, ok := rdx.GetLastVal(data.PlaylistChannelProperty, playlistId); ok && plc != "" {
		playlistChannelTitle = plc
	}

	plvm := &PlaylistViewModel{
		PlaylistId:           playlistId,
		PlaylistClass:        pc,
		NewVideos:            nvc,
		PlaylistTitle:        playlistTitle,
		PlaylistChannelTitle: playlistChannelTitle,
		Watching:             watching,
		Downloading:          downloading,
		SingleFormat:         singleFormat,
	}

	defaultOptions := []VideoOptions{ShowPoster, ShowViewCount, ShowDuration, ShowPublishedDate}

	if videoIds, ok := rdx.GetAllValues(data.PlaylistVideosProperty, playlistId); ok && len(videoIds) > 0 {
		for _, videoId := range videoIds {
			plvm.Videos = append(plvm.Videos, GetVideoViewModel(videoId, rdx, defaultOptions...))
		}
	}
	return plvm
}
