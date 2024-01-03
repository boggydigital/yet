package rest

import (
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/rest/view_models"
	"net/http"
	"slices"
)

const (
	maxPlaylistVideosWatchlist = 3
)

type ListViewModel struct {
	Continue             []*view_models.VideoViewModel
	Watchlist            []*view_models.VideoViewModel
	Downloads            []*view_models.VideoViewModel
	HasNewPlaylistVideos bool
	Playlists            []*view_models.ListPlaylistViewModel
	HasHistory           bool
}

func GetList(w http.ResponseWriter, r *http.Request) {

	// GET /list

	var err error
	rdx, err = rdx.RefreshReader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lvm := &ListViewModel{}

	cwKeys := rdx.Keys(data.VideoProgressProperty)
	if len(cwKeys) > 0 {
		cwKeys, err = rdx.Sort(cwKeys, false, data.VideoTitleProperty)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, id := range cwKeys {
			if ended, ok := rdx.GetFirstVal(data.VideoEndedProperty, id); !ok || ended == "" {
				lvm.Continue = append(lvm.Continue, view_models.GetVideoViewModel(id, rdx,
					view_models.ShowPoster,
					view_models.ShowPublishedDate,
					view_models.ShowRemainingDuration))
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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
			lvm.Watchlist = append(lvm.Watchlist, view_models.GetVideoViewModel(id, rdx,
				view_models.ShowPoster,
				view_models.ShowPublishedDate))
		}

		lvm.HasNewPlaylistVideos = len(newPlaylistVideos) > 0
	}

	plKeys := rdx.Keys(data.PlaylistWatchlistProperty)
	if len(plKeys) > 0 {

		plKeys, err = rdx.Sort(plKeys, false, data.PlaylistTitleProperty, data.PlaylistChannelProperty)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, playlistId := range plKeys {
			lvm.Playlists = append(lvm.Playlists,
				view_models.GetListPlaylistViewModel(playlistId, rdx))
		}
	}

	dqKeys := rdx.Keys(data.VideosDownloadQueueProperty)
	if len(dqKeys) > 0 {

		dqKeys, err = rdx.Sort(dqKeys, false, data.VideoTitleProperty)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, id := range dqKeys {
			lvm.Downloads = append(lvm.Downloads, view_models.GetVideoViewModel(id, rdx,
				view_models.ShowPoster,
				view_models.ShowPublishedDate))
		}
	}

	lvm.HasHistory = len(rdx.Keys(data.VideoEndedProperty)) > 0

	w.Header().Set("Content-Type", "text/html")

	if err := tmpl.ExecuteTemplate(w, "list", lvm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
