package rest

import (
	"encoding/json"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
	"net/http"
)

type EndedRequest struct {
	VideoId string `json:"v"`
}

func PostEnded(w http.ResponseWriter, r *http.Request) {

	// POST /ended
	// {v}

	var err error
	rdx, err = rdx.RefreshWriter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var er EndedRequest
	err = decoder.Decode(&er)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// store completion timestamp
	if err := rdx.ReplaceValues(data.VideoEndedDateProperty, er.VideoId, yeti.FmtNow()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
