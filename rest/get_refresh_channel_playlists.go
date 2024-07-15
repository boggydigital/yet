package rest

import (
	"github.com/boggydigital/yet/yeti"
	"net/http"
)

func GetRefreshChannelPlaylists(w http.ResponseWriter, r *http.Request) {

	// GET /refresh_channel_playlists?id

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

	if err := yeti.GetChannelPlaylistsMetadata(nil, channelId, rdx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/channel_playlists?id="+channelId, http.StatusTemporaryRedirect)

}
