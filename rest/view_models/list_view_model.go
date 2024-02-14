package view_models

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/yet/data"
	"math/rand"
	"slices"
)

type ListViewModel struct {
	Continue             []*VideoViewModel
	Random               *VideoViewModel
	Videos               []*VideoViewModel
	Downloads            []*VideoViewModel
	HasNewPlaylistVideos bool
	PlaylistsOrder       []string
	Playlists            map[string][]*PlaylistViewModel
	HasHistory           bool
}

const (
	playlistsNewVideos   = "Unwatched"
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
			if ended, ok := rdx.GetLastVal(data.VideoEndedProperty, id); !ok || ended == "" {
				lvm.Continue = append(lvm.Continue, GetVideoViewModel(id, rdx,
					ShowPoster,
					ShowPublishedDate,
					ShowDuration,
					ShowProgress))
			}
		}
	}

	if len(lvm.Continue) == 0 {
		// add random video to suggest watching
		pool := make([]string, 0)
		for _, id := range rdx.Keys(data.PlaylistNewVideosProperty) {
			if pnv, ok := rdx.GetLastVal(data.PlaylistNewVideosProperty, id); ok {
				pool = append(pool, pnv)
			}
		}

		if len(pool) > 0 {
			lvm.Random = GetVideoViewModel(pool[rand.Intn(len(pool))], rdx, ShowPoster, ShowPublishedDate, ShowDuration)
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
			if le, ok := rdx.GetLastVal(data.VideoEndedProperty, id); ok && le != "" {
				continue
			}
			if ct, ok := rdx.GetLastVal(data.VideoProgressProperty, id); ok || ct != "" {
				continue
			}
			lvm.Videos = append(lvm.Videos, GetVideoViewModel(id, rdx,
				ShowPoster,
				ShowPublishedDate,
				ShowDuration))
		}

		lvm.HasNewPlaylistVideos = len(newPlaylistVideos) > 0
	}

	plKeys := rdx.Keys(data.PlaylistWatchlistProperty)
	if len(plKeys) > 0 {

		plNewVideos, plNoNewVideos := make([]string, 0, len(plKeys)), make([]string, 0, len(plKeys))

		for _, playlistId := range plKeys {
			if newVideos, ok := rdx.GetAllValues(data.PlaylistNewVideosProperty, playlistId); ok && len(newVideos) > 0 {
				plNewVideos = append(plNewVideos, playlistId)
			} else {
				plNoNewVideos = append(plNoNewVideos, playlistId)
			}
		}

		plNewVideos, err = rdx.Sort(plNewVideos, false, data.PlaylistTitleProperty, data.PlaylistChannelProperty)
		if err != nil {
			return nil, err
		}

		plNoNewVideos, err = rdx.Sort(plNoNewVideos, false, data.PlaylistTitleProperty, data.PlaylistChannelProperty)
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

	dqKeys := rdx.Keys(data.VideosDownloadQueueProperty)
	if len(dqKeys) > 0 {

		dqKeys, err = rdx.Sort(dqKeys, false, data.VideoTitleProperty)
		if err != nil {
			return nil, err
		}

		for _, id := range dqKeys {
			lvm.Downloads = append(lvm.Downloads, GetVideoViewModel(id, rdx,
				ShowPoster,
				ShowDuration,
				ShowPublishedDate))
		}
	}

	lvm.HasHistory = len(rdx.Keys(data.VideoEndedProperty)) > 0

	return lvm, nil
}
