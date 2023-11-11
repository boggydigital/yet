package rest

import (
	"github.com/boggydigital/yet/paths"
	"net/http"
	"os"
)

func GetPoster(w http.ResponseWriter, r *http.Request) {

	// GET /poster?v&q

	videoId := r.URL.Query().Get("v")
	quality := r.URL.Query().Get("q")

	absPosterFilename, err := paths.AbsPosterPath(videoId, quality)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := os.Stat(absPosterFilename); err == nil {
		http.ServeFile(w, r, absPosterFilename)
	} else {
		http.NotFound(w, r)
		return
	}

}
