package rest

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yet/yeti"
	"net/http"
)

func GetRefresh(w http.ResponseWriter, r *http.Request) {

	// GET /refresh?id

	id := r.URL.Query().Get("list")

	if id == "" {
		http.Redirect(w, r, "/new", http.StatusPermanentRedirect)
		return
	}

	metadataDir, err := paths.GetAbsDir(paths.Metadata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	plRdx, err := kvas.NewReduxWriter(metadataDir,
		data.PlaylistTitleProperty,
		data.PlaylistChannelProperty,
		data.PlaylistVideosProperty,
		data.VideoTitleProperty,
		data.VideoOwnerChannelNameProperty)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := yeti.GetPlaylistPageMetadata(nil, id, false, plRdx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/playlist?list="+id, http.StatusTemporaryRedirect)
}
