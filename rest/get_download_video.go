package rest

import (
	"net/http"
	"path"

	"github.com/boggydigital/yet/yeti"
)

func GetDownloadVideo(w http.ResponseWriter, r *http.Request) {

	// Get /download_video/{videoId}

	var err error
	rdx, err = rdx.RefreshWriter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	videoId := r.PathValue("videoId")

	if err = yeti.DownloadVideoMetadataPoster(nil, videoId, yeti.DefaultVideoOptions(), rdx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, path.Join("/watch", videoId), http.StatusTemporaryRedirect)
}
