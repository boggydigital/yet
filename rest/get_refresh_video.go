package rest

import (
	"net/http"
	"path"

	"github.com/boggydigital/yet/yeti"
)

func GetRefreshVideo(w http.ResponseWriter, r *http.Request) {

	// GET /refresh_video/{videoId}

	var err error
	rdx, err = rdx.RefreshWriter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	videoId := r.PathValue("videoId")

	if videoId == "" {
		http.Redirect(w, r, "/list", http.StatusPermanentRedirect)
		return
	}

	videoPage, err := yeti.GetVideoPage(videoId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	videoMetadata := yeti.ExtractMetadata(videoPage)

	for p, values := range videoMetadata {
		if err = rdx.ReplaceValues(p, videoId, values...); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, path.Join("/watch", videoId), http.StatusTemporaryRedirect)
}
