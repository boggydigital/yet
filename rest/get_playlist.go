package rest

import (
	"github.com/boggydigital/yet/data"
	"net/http"
)

const (
	showImagesLimit = 20
)

type PlaylistViewModel struct {
	PlaylistTitle   string
	PlaylistId      string
	AutoDownloading bool
	Videos          []*VideoViewModel
}

func GetPlaylist(w http.ResponseWriter, r *http.Request) {

	// GET /playlist?id

	var err error
	rdx, err = rdx.RefreshReader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id := r.URL.Query().Get("list")

	if id == "" {
		http.Redirect(w, r, "/new", http.StatusPermanentRedirect)
		return
	}

	pt := playlistTitle(id, rdx)

	w.Header().Set("Content-Type", "text/html")

	// playlist specific styles

	plvm := &PlaylistViewModel{
		PlaylistId:    id,
		PlaylistTitle: pt,
	}

	if pdq, ok := rdx.GetFirstVal(data.PlaylistDownloadQueueProperty, id); ok && pdq == data.TrueValue {
		plvm.AutoDownloading = true
	}

	if videoIds, ok := rdx.GetAllValues(data.PlaylistVideosProperty, id); ok && len(videoIds) > 0 {
		for i, videoId := range videoIds {
			var options []VideoOptions
			if i+1 < showImagesLimit {
				options = []VideoOptions{ShowPoster, ShowPublishedDate}
			} else {
				options = []VideoOptions{ShowPublishedDate}
			}
			plvm.Videos = append(plvm.Videos, videoViewModel(videoId, rdx, options...))
		}
	}

	if err := tmpl.ExecuteTemplate(w, "playlist", plvm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
