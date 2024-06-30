package view_models

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
)

type PlaylistViewModel struct {
	PlaylistId           string
	PlaylistTitle        string
	PlaylistChannelTitle string
	BadgeCount           int
	AutoRefresh          bool
	AutoDownload         bool
	DownloadPolicy       data.PlaylistDownloadPolicy
	AllDownloadPolicies  []data.PlaylistDownloadPolicy
	Expand               bool
	PreferSingleFormat   bool
	Videos               []*VideoViewModel
}

func GetPlaylistViewModel(playlistId string, rdx kvas.ReadableRedux) *PlaylistViewModel {

	if playlistId == "" {
		return nil
	}

	pnev := yeti.PlaylistNotEndedVideos(playlistId, rdx)
	badgeCount := len(pnev)

	autoRefresh := false
	if par, ok := rdx.GetLastVal(data.PlaylistAutoRefreshProperty, playlistId); ok && par == data.TrueValue {
		autoRefresh = true
	}

	autoDownload := false
	if pad, ok := rdx.GetLastVal(data.PlaylistAutoDownloadProperty, playlistId); ok && pad == data.TrueValue {
		autoDownload = true
	}

	downloadPolicy := data.DefaultDownloadPolicy
	if pdp, ok := rdx.GetLastVal(data.PlaylistDownloadPolicyProperty, playlistId); ok {
		downloadPolicy = data.ParsePlaylistDownloadPolicy(pdp)
	}

	expand := false
	if pe, ok := rdx.GetLastVal(data.PlaylistExpandProperty, playlistId); ok && pe == data.TrueValue {
		expand = true
	}

	preferSingleFormat := false
	if psf, ok := rdx.GetLastVal(data.PlaylistPreferSingleFormatProperty, playlistId); ok && psf == data.TrueValue {
		preferSingleFormat = true
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
		BadgeCount:           badgeCount,
		PlaylistTitle:        playlistTitle,
		PlaylistChannelTitle: playlistChannelTitle,
		AutoRefresh:          autoRefresh,
		AutoDownload:         autoDownload,
		DownloadPolicy:       downloadPolicy,
		AllDownloadPolicies:  data.AllPlaylistDownloadPolicies(),
		Expand:               expand,
		PreferSingleFormat:   preferSingleFormat,
	}

	defaultOptions := []VideoOptions{ShowPoster, ShowViewCount, ShowDuration, ShowPublishedDate}

	if videoIds, ok := rdx.GetAllValues(data.PlaylistVideosProperty, playlistId); ok && len(videoIds) > 0 {
		for _, videoId := range videoIds {
			plvm.Videos = append(plvm.Videos, GetVideoViewModel(videoId, rdx, defaultOptions...))
		}
	}
	return plvm
}
