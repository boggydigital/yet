package view_models

import (
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
)

type ChannelViewModel struct {
	ChannelId             string
	ChannelTitle          string
	ChannelDescription    string
	ChannelBadgeCount     int
	Videos                []*VideoViewModel
	PlaylistsOrder        []string
	Playlists             map[string]string
	ChannelAutoRefresh    bool
	ChannelAutoDownload   bool
	ChannelDownloadPolicy data.DownloadPolicy
	AllDownloadPolicies   []data.DownloadPolicy
	ChannelExpand         bool
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

	cnev := yeti.ChannelNotEndedVideos(channelId, rdx)
	badgeCount := len(cnev)

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

	expand := false
	if ce, ok := rdx.GetLastVal(data.ChannelExpandProperty, channelId); ok && ce == data.TrueValue {
		expand = true
	}

	var playlistsOrder []string
	if chpls, ok := rdx.GetAllValues(data.ChannelPlaylistsProperty, channelId); ok && len(chpls) > 0 {
		playlistsOrder = chpls
	}

	playlists := make(map[string]string, len(playlistsOrder))
	for _, playlistId := range playlistsOrder {
		if plt, ok := rdx.GetLastVal(data.PlaylistTitleProperty, playlistId); ok && plt != "" {
			playlists[playlistId] = plt
		}
	}

	return &ChannelViewModel{
		ChannelId:             channelId,
		ChannelTitle:          channelTitle,
		ChannelDescription:    channelDescription,
		ChannelBadgeCount:     badgeCount,
		Videos:                channelVideos,
		ChannelAutoRefresh:    autoRefresh,
		ChannelAutoDownload:   autoDownload,
		ChannelDownloadPolicy: downloadPolicy,
		AllDownloadPolicies:   data.AllDownloadPolicies(),
		ChannelExpand:         expand,
		PlaylistsOrder:        playlistsOrder,
		Playlists:             playlists,
	}
}
