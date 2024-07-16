package rest

import (
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/rest/view_models"
	"net/http"
)

func GetPlaylist(w http.ResponseWriter, r *http.Request) {

	// GET /playlist?list

	var err error
	rdx, err = rdx.RefreshWriter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	playlistId := r.URL.Query().Get("list")

	if playlistId == "" {
		http.Redirect(w, r, "/list", http.StatusPermanentRedirect)
		return
	}

	// check if the playlist has no videos and refresh automatically
	if videos, ok := rdx.GetAllValues(data.PlaylistVideosProperty, playlistId); !ok || len(videos) == 0 {
		url := "/refresh_playlist?list=" + playlistId
		http.Redirect(w, r, url, http.StatusPermanentRedirect)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	if err := tmpl.ExecuteTemplate(w, "playlists", view_models.GetPlaylistViewModel(playlistId, rdx)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
