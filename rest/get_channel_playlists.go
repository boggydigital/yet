package rest

import (
	"github.com/boggydigital/yet/rest/view_models"
	"net/http"
)

func GetChannelPlaylists(w http.ResponseWriter, r *http.Request) {

	// GET /channel_playlists?id

	var err error
	rdx, err = rdx.RefreshWriter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	channelId := r.URL.Query().Get("id")

	if channelId == "" {
		http.Redirect(w, r, "/list", http.StatusPermanentRedirect)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	if err := tmpl.ExecuteTemplate(w, "channel_playlists", view_models.GetChannelViewModel(channelId, rdx)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
