package rest

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/yet/data"
	"net/http"
	"path"
	"slices"
	"strings"
)

const (
	maxPlaylistVideosWatchlist = 3
)

type VideoViewModel struct {
	VideoId           string
	VideoUrl          string
	VideoTitle        string
	Class             string
	ShowPoster        bool
	ShowPublishedDate bool
	PublishedDate     string
	ShowEndedDate     bool
	EndedDate         string
}

type PlaylistViewModel struct {
	PlaylistId    string
	PlaylistTitle string
	Class         string
	NewVideos     int
}

type ListViewModel struct {
	Continue             []*VideoViewModel
	Watchlist            []*VideoViewModel
	Downloads            []*VideoViewModel
	HasNewPlaylistVideos bool
	Playlists            []*PlaylistViewModel
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
				lvm.Continue = append(lvm.Continue, videoViewModel(id, rdx, ShowPoster, ShowPublishedDate))
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
			lvm.Watchlist = append(lvm.Watchlist, videoViewModel(id, rdx, ShowPoster, ShowPublishedDate))
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

		for _, id := range plKeys {

			nvc := 0

			if nv, ok := rdx.GetAllValues(data.PlaylistNewVideosProperty, id); ok {
				nvc = len(nv)
			}

			pc := "playlist"
			if nvc == 0 {
				pc += " ended"
			}

			plvm := &PlaylistViewModel{
				PlaylistId:    id,
				PlaylistTitle: playlistTitle(id, rdx),
				Class:         pc,
				NewVideos:     nvc,
			}

			lvm.Playlists = append(lvm.Playlists, plvm)
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
			lvm.Downloads = append(lvm.Downloads, videoViewModel(id, rdx, ShowPoster, ShowPublishedDate))
		}
	}

	lvm.HasHistory = len(rdx.Keys(data.VideoEndedProperty)) > 0

	w.Header().Set("Content-Type", "text/html")

	if err := tmpl.ExecuteTemplate(w, "list", lvm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func playlistTitle(playlistId string, rdx kvas.ReadableRedux) string {
	if plt, ok := rdx.GetFirstVal(data.PlaylistTitleProperty, playlistId); ok && plt != "" {

		if plc, ok := rdx.GetFirstVal(data.PlaylistChannelProperty, playlistId); ok && plc != "" && !strings.Contains(plt, plc) {
			if plt == "Videos" {
				return plc
			} else {
				return path.Join(plc, plt)
			}
		}

		return plt
	}

	return playlistId
}
