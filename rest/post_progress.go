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

	var pr ProgressRequest
	if err := json.NewDecoder(r.Body).Decode(&pr); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data.ProgressMux.Lock()
	data.VideosProgress[pr.VideoId] = []string{trimTime(pr.CurrentTime)}
	data.ProgressMux.Unlock()

}

func trimTime(ts string) string {
	if tt, _, ok := strings.Cut(ts, "."); ok {
		return tt
	}
	return ts
}
