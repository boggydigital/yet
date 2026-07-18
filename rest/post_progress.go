package rest

import (
	"net/http"
	"strings"

	"github.com/boggydigital/yet/data"
)

func PostProgress(w http.ResponseWriter, r *http.Request) {

	// POST /progress/{videoId}/{time}

	var err error
	rdx, err = rdx.RefreshWriter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	videoId := r.PathValue("videoId")
	currentTime := r.PathValue("currentTime")

	if err = rdx.ReplaceValues(data.VideoProgressProperty, videoId, trimTime(currentTime)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func trimTime(ts string) string {
	if tt, _, ok := strings.Cut(ts, "."); ok {
		return tt
	}
	return ts
}
