package view_models

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/yet/data"
	"slices"
)

type ListViewModel struct {
	Continue             []*VideoViewModel
	Videos               []*VideoViewModel
	Downloads            []*VideoViewModel
	HasNewPlaylistVideos bool
	Playlists            []*PlaylistViewModel
	HasHistory           bool
}

func GetListViewModel(rdx kvas.ReadableRedux) (*ListViewModel, error) {
	lvm := &ListViewModel{}

	var err error

	cwKeys := rdx.Keys(data.VideoProgressProperty)
	if len(cwKeys) > 0 {
		cwKeys, err = rdx.Sort(cwKeys, false, data.VideoTitleProperty)
		if err != nil {
			return nil, err
		}
		for _, id := range cwKeys {
			if ended, ok := rdx.GetFirstVal(data.VideoEndedProperty, id); !ok || ended == "" {
				lvm.Continue = append(lvm.Continue, GetVideoViewModel(id, rdx,
					ShowPoster,
					ShowPublishedDate,
					ShowRemainingDuration))
			}
		}
	}

	pldq := rdx.Keys(data.PlaylistDownloadQueueProperty)
	newPlaylistVideos := make([]string, 0, len(pldq))

	for _, pl := range pldq {
		if nv, ok := rdx.GetAllValues(data.PlaylistNewVideosProperty, pl); ok {
			newPlaylistVideos = append(newPlaylistVideos, nv...)
		}
	}

	wlKeys := rdx.Keys(data.VideosWatchlistProperty)
	if len(wlKeys) > 0 {

		wlKeys, err = rdx.Sort(wlKeys, false, data.VideoTitleProperty)
		if err != nil {
			return nil, err
		}

		for _, id := range wlKeys {
			if slices.Contains(newPlaylistVideos, id) {
				continue
			}
			if le, ok := rdx.GetFirstVal(data.VideoEndedProperty, id); ok && le != "" {
				continue
			}
			if ct, ok := rdx.GetFirstVal(data.VideoProgressProperty, id); ok || ct != "" {
				continue
			}
			lvm.Videos = append(lvm.Videos, GetVideoViewModel(id, rdx,
				ShowPoster,
				ShowPublishedDate))
		}

		lvm.HasNewPlaylistVideos = len(newPlaylistVideos) > 0
	}

	plKeys := rdx.Keys(data.PlaylistWatchlistProperty)
	if len(plKeys) > 0 {

		plKeys, err = rdx.Sort(plKeys, false, data.PlaylistTitleProperty, data.PlaylistChannelProperty)
		if err != nil {
			return nil, err
		}

		for _, playlistId := range plKeys {
			lvm.Playlists = append(lvm.Playlists, GetPlaylistViewModel(playlistId, rdx))
		}
	}

	dqKeys := rdx.Keys(data.VideosDownloadQueueProperty)
	if len(dqKeys) > 0 {

		dqKeys, err = rdx.Sort(dqKeys, false, data.VideoTitleProperty)
		if err != nil {
			return nil, err
		}

		for _, id := range dqKeys {
			lvm.Downloads = append(lvm.Downloads, GetVideoViewModel(id, rdx,
				ShowPoster,
				ShowPublishedDate))
		}
	}

	lvm.HasHistory = len(rdx.Keys(data.VideoEndedProperty)) > 0

	return lvm, nil
}
