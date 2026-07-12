package rest

import (
	"net/http"
	"path"

	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
)

type QueueDownloadRequest struct {
	VideoId string `json:"v"`
}

func GetQueueDownload(w http.ResponseWriter, r *http.Request) {

	// Get /queue_download/{videoId}

	var err error
	rdx, err = rdx.RefreshWriter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	videoId := r.PathValue("videoId")

	// store completion timestamp
	if err = rdx.AddValues(data.VideoDownloadQueuedProperty, videoId, yeti.FmtNow()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, path.Join("/watch", videoId), http.StatusTemporaryRedirect)
}
