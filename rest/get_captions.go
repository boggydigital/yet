package rest

import (
	"github.com/boggydigital/yet/data"
	"net/http"
	"os"
)

func GetCaptions(w http.ResponseWriter, r *http.Request) {

	// GET /captions?v&l

	videoId := r.URL.Query().Get("v")
	lang := r.URL.Query().Get("l")

	absCaptionsFilename, err := data.AbsCaptionsTrackPath(videoId, lang)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := os.Stat(absCaptionsFilename); err == nil {
		http.ServeFile(w, r, absCaptionsFilename)
	} else {
		http.NotFound(w, r)
	}

}
