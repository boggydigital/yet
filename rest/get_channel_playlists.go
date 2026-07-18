package rest

import (
	"net/http"
	"path"

	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/rest/view_models"
)

func GetChannelPlaylists(w http.ResponseWriter, r *http.Request) {

	// GET /channel_playlists/{channelId}

	var err error
	rdx, err = rdx.RefreshWriter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	channelId := r.PathValue("channelId")

	if channelId == "" {
		http.Redirect(w, r, "/list", http.StatusPermanentRedirect)
		return
	}

	// check if the channel has no playlists and refresh automatically
	if playlists, ok := rdx.GetAllValues(data.ChannelPlaylistsProperty, channelId); !ok || len(playlists) == 0 {
		http.Redirect(w, r, path.Join("/refresh_channel_playlists", channelId), http.StatusPermanentRedirect)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	if err = tmpl.ExecuteTemplate(w, "channel_playlists", view_models.GetChannelViewModel(channelId, rdx)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
