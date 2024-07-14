package view_models

import (
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/yet/data"
)

type ChannelViewModel struct {
	ChannelId          string
	ChannelTitle       string
	ChannelDescription string
	Videos             []*VideoViewModel
	//PlaylistsOrder     []string
	//Playlists          map[string]string
	//PlaylistsVideos    map[string][]*VideoViewModel
}

func GetChannelViewModel(channelId string, rdx kevlar.ReadableRedux) *ChannelViewModel {
	channelTitle := channelId
	if ct, ok := rdx.GetLastVal(data.ChannelTitleProperty, channelId); ok && ct != "" {
		channelTitle = ct
	}

	channelDescription := ""
	if cd, ok := rdx.GetLastVal(data.ChannelDescriptionProperty, channelId); ok && cd != "" {
		channelDescription = cd
	}

	var channelVideos []*VideoViewModel
	if chvs, ok := rdx.GetAllValues(data.ChannelVideosProperty, channelId); ok && len(chvs) > 0 {
		for _, videoId := range chvs {
			channelVideos = append(channelVideos, GetVideoViewModel(videoId, rdx, ShowPoster, ShowPublishedDate, ShowViewCount))
		}
	}

	//var playlistsOrder []string
	//if chpls, ok := rdx.GetAllValues(data.ChannelPlaylistsProperty, channelId); ok && len(chpls) > 0 {
	//	playlistsOrder = chpls
	//}
	//
	//playlists := make(map[string]string, len(playlistsOrder))
	//for _, playlistId := range playlistsOrder {
	//	if plt, ok := rdx.GetLastVal(data.PlaylistTitleProperty, playlistId); ok && plt != "" {
	//		playlists[playlistId] = plt
	//	}
	//}
	//
	//playlistVideos := make(map[string][]*VideoViewModel, len(playlistsOrder))
	//for _, playlistId := range playlistsOrder {
	//	if plvs, ok := rdx.GetAllValues(data.PlaylistVideosProperty, playlistId); ok {
	//		for _, videoId := range plvs {
	//			playlistVideos[playlistId] = append(playlistVideos[playlistId], GetVideoViewModel(videoId, rdx, ShowPublishedDate, ShowViewCount))
	//		}
	//	}
	//}

	return &ChannelViewModel{
		ChannelId:          channelId,
		ChannelTitle:       channelTitle,
		ChannelDescription: channelDescription,
		Videos:             channelVideos,
		//PlaylistsOrder:     playlistsOrder,
		//Playlists:          playlists,
		//PlaylistsVideos:    playlistVideos,
	}
}
