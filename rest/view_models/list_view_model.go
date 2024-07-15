package view_models

import (
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
)

type ListViewModel struct {
	Continue       []*VideoViewModel
	Videos         []*VideoViewModel
	Downloads      []*VideoViewModel
	ChannelsOrder  []string
	Channels       map[string][]*ChannelViewModel
	PlaylistsOrder []string
	Playlists      map[string][]*PlaylistViewModel
	Favorites      []*VideoViewModel
	HasHistory     bool
}

const (
	newItems   = "New"
	noNewItems = "Watched"
)

var (
	channelsOrder  = []string{newItems, noNewItems}
	playlistsOrder = []string{newItems, noNewItems}
)

func GetListViewModel(rdx kevlar.ReadableRedux) (*ListViewModel, error) {
	lvm := &ListViewModel{
		ChannelsOrder:  channelsOrder,
		Channels:       make(map[string][]*ChannelViewModel),
		PlaylistsOrder: playlistsOrder,
		Playlists:      make(map[string][]*PlaylistViewModel),
	}

	if videosContinue, err := getVideosProgress(rdx); err == nil {
		for _, id := range videosContinue {
			lvm.Continue = append(lvm.Continue, GetVideoViewModel(id, rdx,
				ShowPoster,
				ShowPublishedDate,
				ShowDuration,
				ShowProgress))
		}
	} else {
		return nil, err
	}

	if videosDownloads, err := getVideoDownloads(rdx); err == nil {
		for _, id := range videosDownloads {
			lvm.Videos = append(lvm.Videos, GetVideoViewModel(id, rdx,
				ShowPoster,
				ShowPublishedDate,
				ShowDuration))
		}
	} else {
		return nil, err
	}

	// channels
	if channelsVideos, err := getChannelsVideos(rdx); err == nil {
		for _, updates := range []string{newItems, noNewItems} {
			for _, channelId := range channelsVideos[updates] {
				lvm.Channels[updates] = append(lvm.Channels[updates], GetChannelViewModel(channelId, rdx))
			}
		}
	}

	// playlists
	if playlistsVideos, err := getPlaylistsVideos(rdx); err == nil {
		for _, updates := range []string{newItems, noNewItems} {
			for _, playlistId := range playlistsVideos[updates] {
				lvm.Playlists[updates] = append(lvm.Playlists[updates], GetPlaylistViewModel(playlistId, rdx))
			}
		}
	}

	if queuedDownloads, err := getQueuedDownloads(rdx); err == nil {
		for _, id := range queuedDownloads {
			lvm.Downloads = append(lvm.Downloads, GetVideoViewModel(id, rdx,
				ShowPoster,
				ShowDuration,
				ShowPublishedDate))
		}
	} else {
		return nil, err
	}

	if favoriteVideos, err := getFavoriteVideos(rdx); err == nil {
		for _, id := range favoriteVideos {
			lvm.Favorites = append(lvm.Favorites, GetVideoViewModel(id, rdx,
				ShowPoster,
				ShowDuration,
				ShowPublishedDate))
		}
	} else {
		return nil, err
	}

	lvm.HasHistory = len(rdx.Keys(data.VideoEndedDateProperty)) > 0

	return lvm, nil
}

func getVideosProgress(rdx kevlar.ReadableRedux) ([]string, error) {
	vpKeys := rdx.Keys(data.VideoProgressProperty)
	cvs := make([]string, 0)
	var err error

	if len(vpKeys) == 0 {
		return cvs, nil
	}

	for _, id := range vpKeys {
		if et, ok := rdx.GetLastVal(data.VideoEndedDateProperty, id); ok && et != "" {
			continue
		}
		cvs = append(cvs, id)
	}

	if cvs, err = rdx.Sort(cvs, false, data.VideoTitleProperty); err == nil {
		return cvs, nil
	} else {
		return nil, err
	}
}

func getVideoDownloads(rdx kevlar.ReadableRedux) ([]string, error) {

	dcKeys := rdx.Keys(data.VideoDownloadCompletedProperty)
	dvs := make([]string, 0, len(dcKeys))

	if len(dcKeys) == 0 {
		return dvs, nil
	}

	// videos is all downloaded videos that are not:
	// - in history (ended)
	// - in continue (have progress)
	// - is favorite
	// - in any auto-refreshing channel
	// - in any auto-refreshing playlist

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

		// check if this video is an auto-refreshing channel video
		skip := false
		for _, channelId := range rdx.Keys(data.ChannelAutoRefreshProperty) {
			if rdx.HasValue(data.ChannelVideosProperty, channelId, id) {
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		// check if this video is an auto-refreshing playlist video
		skip = false
		for _, playlistId := range rdx.Keys(data.PlaylistAutoRefreshProperty) {
			if rdx.HasValue(data.PlaylistVideosProperty, playlistId, id) {
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		dvs = append(dvs, id)
	}

	var err error
	if dvs, err = rdx.Sort(dvs, false, data.VideoTitleProperty); err == nil {
		return dvs, nil
	} else {
		return nil, err
	}
}

func getChannelsVideos(rdx kevlar.ReadableRedux) (map[string][]string, error) {
	chKeys := rdx.Keys(data.ChannelAutoRefreshProperty)
	chs := make(map[string][]string)

	if len(chKeys) == 0 {
		return chs, nil
	}

	chNewVideos, chNoNewVideos := make([]string, 0, len(chKeys)), make([]string, 0, len(chKeys))

	for _, channelId := range chKeys {
		if newVideos := yeti.ChannelNotEndedVideos(channelId, rdx); len(newVideos) > 0 {
			chNewVideos = append(chNewVideos, channelId)
		} else {
			chNoNewVideos = append(chNoNewVideos, channelId)
		}
	}

	var err error
	chNewVideos, err = rdx.Sort(chNewVideos, false, data.ChannelTitleProperty)
	if err != nil {
		return nil, err
	}

	chNoNewVideos, err = rdx.Sort(chNoNewVideos, false, data.ChannelTitleProperty)
	if err != nil {
		return nil, err
	}

	chs[newItems] = chNewVideos
	chs[noNewItems] = chNoNewVideos

	return chs, nil
}

func getPlaylistsVideos(rdx kevlar.ReadableRedux) (map[string][]string, error) {
	plKeys := rdx.Keys(data.PlaylistAutoRefreshProperty)
	pls := make(map[string][]string)

	if len(plKeys) == 0 {
		return pls, nil
	}

	plNewVideos, plNoNewVideos := make([]string, 0, len(plKeys)), make([]string, 0, len(plKeys))

	for _, playlistId := range plKeys {
		if newVideos := yeti.PlaylistNotEndedVideos(playlistId, rdx); len(newVideos) > 0 {
			plNewVideos = append(plNewVideos, playlistId)
		} else {
			plNoNewVideos = append(plNoNewVideos, playlistId)
		}
	}

	var err error
	plNewVideos, err = rdx.Sort(plNewVideos, false, data.PlaylistTitleProperty, data.PlaylistChannelProperty)
	if err != nil {
		return nil, err
	}

	plNoNewVideos, err = rdx.Sort(plNoNewVideos, false, data.PlaylistTitleProperty, data.PlaylistChannelProperty)
	if err != nil {
		return nil, err
	}

	pls[newItems] = plNewVideos
	pls[noNewItems] = plNoNewVideos

	return pls, nil
}

func getQueuedDownloads(rdx kevlar.ReadableRedux) ([]string, error) {

	dqKeys := rdx.Keys(data.VideoDownloadQueuedProperty)
	qds := make([]string, 0, len(dqKeys))

	if len(dqKeys) == 0 {
		return qds, nil
	}

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

		qds = append(qds, id)
	}

	var err error
	if qds, err = rdx.Sort(qds, false, data.VideoTitleProperty); err == nil {
		return qds, nil
	} else {
		return nil, err
	}
}

func getFavoriteVideos(rdx kevlar.ReadableRedux) ([]string, error) {
	fvs := rdx.Keys(data.VideoFavoriteProperty)
	var err error
	if fvs, err = rdx.Sort(fvs, false, data.VideoTitleProperty); err == nil {
		return fvs, nil
	} else {
		return nil, err
	}
}
