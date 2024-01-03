package view_models

import (
	"fmt"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/yet/data"
	"strings"
)

type ListPlaylistViewModel struct {
	PlaylistId    string
	PlaylistTitle string
	Class         string
	NewVideos     int
}

type PlaylistViewModel struct {
	PlaylistTitle   string
	PlaylistId      string
	AutoDownloading bool
	Videos          []*VideoViewModel
}

func GetListPlaylistViewModel(playlistId string, rdx kvas.ReadableRedux) *ListPlaylistViewModel {
	nvc := 0

	if nv, ok := rdx.GetAllValues(data.PlaylistNewVideosProperty, playlistId); ok {
		nvc = len(nv)
	}

	pc := "playlist"
	if nvc == 0 {
		pc += " ended"
	}

	return &ListPlaylistViewModel{
		PlaylistId:    playlistId,
		PlaylistTitle: PlaylistTitle(playlistId, rdx),
		Class:         pc,
		NewVideos:     nvc,
	}
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
