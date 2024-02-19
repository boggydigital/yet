package rest

import (
	"encoding/json"
	"github.com/boggydigital/yet/data"
	"net/http"
	"strings"
)

type ProgressRequest struct {
	VideoId     string `json:"v"`
	CurrentTime string `json:"t"`
}

func PostProgress(w http.ResponseWriter, r *http.Request) {

	// POST /progress
	// {v, t}

	var err error
	progressRdx, err = progressRdx.RefreshWriter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var pr ProgressRequest
	err = decoder.Decode(&pr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := progressRdx.ReplaceValues(data.VideoProgressProperty, pr.VideoId, trimTime(pr.CurrentTime)); err != nil {
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
