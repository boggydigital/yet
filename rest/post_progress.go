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

func (pr ProgressRequest) TrimTime() string {
	if tt, _, ok := strings.Cut(pr.CurrentTime, "."); ok {
		return tt
	}
	return pr.CurrentTime
}

func PostProgress(w http.ResponseWriter, r *http.Request) {

	// POST /progress
	// {v, t}

	decoder := json.NewDecoder(r.Body)
	var pr ProgressRequest
	err := decoder.Decode(&pr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := rdx.ReplaceValues(data.VideoProgressProperty, pr.VideoId, pr.TrimTime()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
