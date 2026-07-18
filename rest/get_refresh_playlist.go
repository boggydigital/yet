package rest

import (
	"net/http"
	"path"

	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
)

func GetRefreshPlaylist(w http.ResponseWriter, r *http.Request) {

	// GET /refresh_playlist?list

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

	expand := false
	if exp, ok := rdx.GetLastVal(data.PlaylistExpandProperty, playlistId); ok && exp == data.TrueValue {
		expand = true
	}

	if err = yeti.GetPlaylistMetadata(nil, playlistId, expand, rdx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, path.Join("/playlist", playlistId), http.StatusTemporaryRedirect)
}
