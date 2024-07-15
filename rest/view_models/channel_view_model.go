package view_models

import (
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/yet/data"
)

type ChannelViewModel struct {
	ChannelId           string
	ChannelTitle        string
	ChannelDescription  string
	Videos              []*VideoViewModel
	AutoRefresh         bool
	AutoDownload        bool
	DownloadPolicy      data.DownloadPolicy
	AllDownloadPolicies []data.DownloadPolicy
	PreferSingleFormat  bool
	Expand              bool

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
			channelVideos = append(channelVideos, GetVideoViewModel(videoId, rdx, ShowPoster, ShowDuration, ShowPublishedDate, ShowViewCount))
		}
	}

	autoRefresh := false
	if par, ok := rdx.GetLastVal(data.ChannelAutoRefreshProperty, channelId); ok && par == data.TrueValue {
		autoRefresh = true
	}

	autoDownload := false
	if pad, ok := rdx.GetLastVal(data.ChannelAutoDownloadProperty, channelId); ok && pad == data.TrueValue {
		autoDownload = true
	}

	downloadPolicy := data.DefaultDownloadPolicy
	if pdp, ok := rdx.GetLastVal(data.ChannelDownloadPolicyProperty, channelId); ok {
		downloadPolicy = data.ParseDownloadPolicy(pdp)
	}

	preferSingleFormat := false
	if psf, ok := rdx.GetLastVal(data.ChannelPreferSingleFormatProperty, channelId); ok && psf == data.TrueValue {
		preferSingleFormat = true
	}

	expand := false
	if pe, ok := rdx.GetLastVal(data.ChannelExpandProperty, channelId); ok && pe == data.TrueValue {
		expand = true
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
		ChannelId:           channelId,
		ChannelTitle:        channelTitle,
		ChannelDescription:  channelDescription,
		Videos:              channelVideos,
		AutoRefresh:         autoRefresh,
		AutoDownload:        autoDownload,
		DownloadPolicy:      downloadPolicy,
		AllDownloadPolicies: data.AllDownloadPolicies(),
		PreferSingleFormat:  preferSingleFormat,
		Expand:              expand,
		//PlaylistsOrder:     playlistsOrder,
		//Playlists:          playlists,
		//PlaylistsVideos:    playlistVideos,

	}
}
