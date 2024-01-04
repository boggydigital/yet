package view_models

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/yet/data"
)

type ChannelViewModel struct {
	ChannelId          string
	ChannelTitle       string
	ChannelDescription string
	PlaylistsOrder     []string
	Playlists          map[string]string
	PlaylistsVideos    map[string][]*VideoViewModel
}

func GetChannelViewModel(channelId string, rdx kvas.ReadableRedux) *ChannelViewModel {
	channelTitle := channelId
	if ct, ok := rdx.GetFirstVal(data.ChannelTitleProperty, channelId); ok && ct != "" {
		channelTitle = ct
	}

	channelDescription := ""
	if cd, ok := rdx.GetFirstVal(data.ChannelDescriptionProperty, channelId); ok && cd != "" {
		channelDescription = cd
	}

	var playlistsOrder []string
	if chpls, ok := rdx.GetAllValues(data.ChannelPlaylistsProperty, channelId); ok && len(chpls) > 0 {
		playlistsOrder = chpls
	}

	playlists := make(map[string]string, len(playlistsOrder))
	for _, playlistId := range playlistsOrder {
		if plt, ok := rdx.GetFirstVal(data.PlaylistTitleProperty, playlistId); ok && plt != "" {
			playlists[playlistId] = plt
		}
	}

	playlistVideos := make(map[string][]*VideoViewModel, len(playlistsOrder))
	for _, playlistId := range playlistsOrder {
		if plvs, ok := rdx.GetAllValues(data.PlaylistVideosProperty, playlistId); ok {
			for _, videoId := range plvs {
				playlistVideos[playlistId] = append(playlistVideos[playlistId], GetVideoViewModel(videoId, rdx, ShowPublishedDate, ShowViewCount))
			}
		}
	}

	return &ChannelViewModel{
		ChannelId:          channelId,
		ChannelTitle:       channelTitle,
		ChannelDescription: channelDescription,
		PlaylistsOrder:     playlistsOrder,
		Playlists:          playlists,
		PlaylistsVideos:    playlistVideos,
	}
}
