package rest

import (
	"encoding/json"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
	"net/http"
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

	decoder := json.NewDecoder(r.Body)
	var qdr QueueDownloadRequest
	err = decoder.Decode(&qdr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// store completion timestamp
	if err := rdx.AddValues(data.VideoDownloadQueuedProperty, qdr.VideoId, yeti.FmtNow()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
