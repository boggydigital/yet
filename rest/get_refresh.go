package rest

import (
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

	if err := yeti.GetPlaylistPageMetadata(nil, id, false, rdx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/playlist?list="+id, http.StatusTemporaryRedirect)
}
