package rest

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/pathways"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yet/yeti"
	"net/http"
)

func GetRefreshPlaylist(w http.ResponseWriter, r *http.Request) {

	// GET /refresh_playlist?list

	playlistId := r.URL.Query().Get("list")

	if playlistId == "" {
		http.Redirect(w, r, "/list", http.StatusPermanentRedirect)
		return
	}

	metadataDir, err := pathways.GetAbsDir(paths.Metadata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	plRdx, err := kvas.NewReduxWriter(metadataDir, data.AllProperties()...)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := yeti.GetPlaylistMetadata(nil, playlistId, false, plRdx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/playlist?list="+playlistId, http.StatusTemporaryRedirect)
}
