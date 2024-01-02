package rest

import (
	"encoding/json"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"net/http"
	"strings"
)

type ProgressRequest struct {
	VideoId     string `json:"v"`
	CurrentTime string `json:"t"`
	Duration    string `json:"d"`
}

func PostProgress(w http.ResponseWriter, r *http.Request) {

	// POST /progress
	// {v, t}

	metadataDir, err := paths.GetAbsDir(paths.Metadata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	progRdx, err := kvas.NewReduxWriter(metadataDir,
		data.VideoProgressProperty,
		data.VideoDurationProperty)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var pr ProgressRequest
	err = decoder.Decode(&pr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := progRdx.ReplaceValues(data.VideoProgressProperty, pr.VideoId, trimTime(pr.CurrentTime)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	td := trimTime(pr.Duration)
	if !progRdx.HasValue(data.VideoDurationProperty, pr.VideoId, td) {
		if err := progRdx.ReplaceValues(data.VideoDurationProperty, pr.VideoId, td); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func trimTime(ts string) string {
	if tt, _, ok := strings.Cut(ts, "."); ok {
		return tt
	}
	return ts
}
