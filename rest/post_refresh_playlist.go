package rest

import (
	"net/http"
	"path"

	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
)

func PostRefreshPlaylist(w http.ResponseWriter, r *http.Request) {

	// POST /refresh_playlist/{playlistId}

	var err error
	rdx, err = rdx.RefreshWriter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	playlistId := r.PathValue("playlistId")

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

	if rdx.HasKey(data.PlaylistAutoDownloadProperty, playlistId) {
		if err = yeti.QueuePlaylistDownloads(rdx, playlistId); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Location", path.Join("/playlist", playlistId))
	w.WriteHeader(http.StatusSeeOther)
	return
}
