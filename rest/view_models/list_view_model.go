package view_models

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
)

type ListViewModel struct {
	Continue       []*VideoViewModel
	Videos         []*VideoViewModel
	Downloads      []*VideoViewModel
	PlaylistsOrder []string
	Playlists      map[string][]*PlaylistViewModel
	HasHistory     bool
}

const (
	playlistsNewVideos   = "New"
	playlistsNoNewVideos = "Watched"
)

var (
	playlistsOrder = []string{playlistsNewVideos, playlistsNoNewVideos}
)

func GetListViewModel(rdx kvas.ReadableRedux) (*ListViewModel, error) {
	lvm := &ListViewModel{
		PlaylistsOrder: playlistsOrder,
		Playlists:      make(map[string][]*PlaylistViewModel),
	}

	var err error

	cwKeys := rdx.Keys(data.VideoProgressProperty)
	if len(cwKeys) > 0 {
		cwKeys, err = rdx.Sort(cwKeys, false, data.VideoTitleProperty)
		if err != nil {
			return nil, err
		}
		for _, id := range cwKeys {
			if et, ok := rdx.GetLastVal(data.VideoEndedDateProperty, id); !ok || et == "" {
				lvm.Continue = append(lvm.Continue, GetVideoViewModel(id, rdx,
					ShowPoster,
					ShowPublishedDate,
					ShowDuration,
					ShowProgress))
			}
		}
	}

	//pldq := rdx.Keys(data.PlaylistAutoDownloadProperty)
	//newPlaylistVideos := make([]string, 0, len(pldq))

	//for _, pl := range pldq {
	//	if nv, ok := rdx.GetAllValues(data.PlaylistNewVideosProperty, pl); ok {
	//		newPlaylistVideos = append(newPlaylistVideos, nv...)
	//	}
	//}

	//wlKeys := rdx.Keys(data.VideosWatchlistProperty)
	//if len(wlKeys) > 0 {
	//
	//	wlKeys, err = rdx.Sort(wlKeys, false, data.VideoTitleProperty)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	for _, id := range wlKeys {
	//		if slices.Contains(newPlaylistVideos, id) {
	//			continue
	//		}
	//		if le, ok := rdx.GetLastVal(data.VideoEndedProperty, id); ok && le != "" {
	//			continue
	//		}
	//		if ct, ok := rdx.GetLastVal(data.VideoProgressProperty, id); ok || ct != "" {
	//			continue
	//		}
	//		lvm.Videos = append(lvm.Videos, GetVideoViewModel(id, rdx,
	//			ShowPoster,
	//			ShowPublishedDate,
	//			ShowDuration))
	//	}
	//}

	plKeys := rdx.Keys(data.PlaylistAutoRefreshProperty)
	if len(plKeys) > 0 {

		plNewVideos, plNoNewVideos := make([]string, 0, len(plKeys)), make([]string, 0, len(plKeys))

		for _, playlistId := range plKeys {
			if newVideos := yeti.PlaylistNotEndedVideos(playlistId, rdx); len(newVideos) > 0 {
				plNewVideos = append(plNewVideos, playlistId)
			} else {
				plNoNewVideos = append(plNoNewVideos, playlistId)
			}
		}

		plNewVideos, err = rdx.Sort(plNewVideos, false, data.PlaylistChannelProperty, data.PlaylistTitleProperty)
		if err != nil {
			return nil, err
		}

		plNoNewVideos, err = rdx.Sort(plNoNewVideos, false, data.PlaylistChannelProperty, data.PlaylistTitleProperty)
		if err != nil {
			return nil, err
		}

		for _, playlistId := range plNewVideos {
			lvm.Playlists[playlistsNewVideos] = append(lvm.Playlists[playlistsNewVideos], GetPlaylistViewModel(playlistId, rdx))
		}

		for _, playlistId := range plNoNewVideos {
			lvm.Playlists[playlistsNoNewVideos] = append(lvm.Playlists[playlistsNoNewVideos], GetPlaylistViewModel(playlistId, rdx))
		}
	}

	dqKeys := rdx.Keys(data.VideoDownloadQueuedProperty)
	if len(dqKeys) > 0 {

		dqKeys, err = rdx.Sort(dqKeys, false, data.VideoTitleProperty)
		if err != nil {
			return nil, err
		}

		for _, id := range dqKeys {
			if dcd, ok := rdx.GetLastVal(data.VideoDownloadCompletedProperty, id); ok && dcd != "" {
				continue
			}
			lvm.Downloads = append(lvm.Downloads, GetVideoViewModel(id, rdx,
				ShowPoster,
				ShowDuration,
				ShowPublishedDate))
		}
	}

	lvm.HasHistory = len(rdx.Keys(data.VideoEndedDateProperty)) > 0

	return lvm, nil
}
