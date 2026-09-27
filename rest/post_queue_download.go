package rest

import (
	"net/http"
	"path"

	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
)

func PostQueueDownload(w http.ResponseWriter, r *http.Request) {

	// POST /queue_download/{videoId}

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

	w.Header().Set("Location", path.Join("/watch", videoId))
	w.WriteHeader(http.StatusSeeOther)
	return
}
