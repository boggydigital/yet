package rest

import (
	"encoding/json"
	"github.com/boggydigital/yet/data"
	"net/http"
	"time"
)

type EndedRequest struct {
	VideoId string `json:"v"`
}

func PostEnded(w http.ResponseWriter, r *http.Request) {

	// POST /ended
	// {v}

	decoder := json.NewDecoder(r.Body)
	var er EndedRequest
	err := decoder.Decode(&er)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	currentTime := time.Now().Format(http.TimeFormat)
	if err := eprdx.ReplaceValues(data.VideoEndedProperty, er.VideoId, currentTime); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
