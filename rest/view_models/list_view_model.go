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
	Favorites      []*VideoViewModel
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

	// videos is all downloaded videos that are not:
	// - in history (ended)
	// - in continue (have progress)
	// - is favorite
	// - in any auto-refreshing playlist
	dcKeys := rdx.Keys(data.VideoDownloadCompletedProperty)
	if len(dcKeys) > 0 {

		notPlaylistDcKeys := make([]string, 0, len(dcKeys))

		for _, id := range dcKeys {

			if rdx.HasKey(data.VideoEndedDateProperty, id) {
				continue
			}
			if rdx.HasKey(data.VideoProgressProperty, id) {
				continue
			}
			if rdx.HasKey(data.VideoFavoriteProperty, id) {
				continue
			}

			skip := false
			for _, playlistId := range rdx.Keys(data.PlaylistAutoRefreshProperty) {
				if rdx.HasValue(data.PlaylistVideosProperty, playlistId, id) {
					skip = true
					break
				}
			}
			if skip {
				continue
			}

			notPlaylistDcKeys = append(notPlaylistDcKeys, id)
		}

		notPlaylistDcKeys, err = rdx.Sort(notPlaylistDcKeys, false, data.VideoTitleProperty)
		if err != nil {
			return nil, err
		}

		for _, id := range notPlaylistDcKeys {
			lvm.Videos = append(lvm.Videos, GetVideoViewModel(id, rdx,
				ShowPoster,
				ShowPublishedDate,
				ShowDuration))
		}
	}

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

		activeDqKeys := make([]string, 0, len(dqKeys))

		for _, id := range dqKeys {

			dqTime := ""
			if dqt, ok := rdx.GetLastVal(data.VideoDownloadQueuedProperty, id); ok {
				dqTime = dqt
			}

			// only continue if download was completed _after_ it was queued,
			// meaning it wasn't re-queued again after completion
			if dcd, ok := rdx.GetLastVal(data.VideoDownloadCompletedProperty, id); ok && dcd > dqTime {
				continue
			}
			activeDqKeys = append(activeDqKeys, id)
		}

		activeDqKeys, err = rdx.Sort(activeDqKeys, false, data.VideoTitleProperty)
		if err != nil {
			return nil, err
		}

		for _, id := range activeDqKeys {
			lvm.Downloads = append(lvm.Downloads, GetVideoViewModel(id, rdx,
				ShowPoster,
				ShowDuration,
				ShowPublishedDate))
		}
	}

	fvKeys := rdx.Keys(data.VideoFavoriteProperty)
	fvKeys, err = rdx.Sort(fvKeys, false, data.VideoTitleProperty)
	if err != nil {
		return nil, err
	}

	for _, id := range fvKeys {
		lvm.Favorites = append(lvm.Favorites, GetVideoViewModel(id, rdx,
			ShowPoster,
			ShowDuration,
			ShowPublishedDate))
	}

	lvm.HasHistory = len(rdx.Keys(data.VideoEndedDateProperty)) > 0

	return lvm, nil
}
