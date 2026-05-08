package rest

import (
	"encoding/json/v2"
	"net/http"

	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
)

type QueueDownloadRequest struct {
	VideoId string `json:"v"`
}

func PostQueueDownload(w http.ResponseWriter, r *http.Request) {

	// POST /queue_download
	// {v}

	var err error
	rdx, err = rdx.RefreshWriter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var qdr QueueDownloadRequest

	if err = json.UnmarshalRead(r.Body, &qdr); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// store completion timestamp
	if err = rdx.AddValues(data.VideoDownloadQueuedProperty, qdr.VideoId, yeti.FmtNow()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
