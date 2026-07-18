package rest

import (
	"net/http"
	"path"

	"github.com/boggydigital/yet/yeti"
)

func GetRefreshChannelPlaylists(w http.ResponseWriter, r *http.Request) {

	// GET /refresh_channel_playlists/{channelId}

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

	if err = yeti.GetChannelPlaylistsMetadata(nil, channelId, rdx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, path.Join("/channel_playlists", channelId), http.StatusTemporaryRedirect)

}
