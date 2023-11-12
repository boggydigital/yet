package rest

import (
	"encoding/json"
	"github.com/boggydigital/yet/data"
	"net/http"
)

type ProgressRequest struct {
	VideoId     string `json:"videoId"`
	CurrentTime string `json:"currentTime"`
}

func PostProgress(w http.ResponseWriter, r *http.Request) {

	// POST /progress
	// {videoId, currentTime}

	decoder := json.NewDecoder(r.Body)
	var pr ProgressRequest
	err := decoder.Decode(&pr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := epRxa.ReplaceValues(data.VideoProgressProperty, pr.VideoId, pr.CurrentTime); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
