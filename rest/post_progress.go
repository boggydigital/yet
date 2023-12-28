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

	metadataDir, err := paths.GetAbsDir(paths.Metadata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	progRdx, err := kvas.NewReduxWriter(metadataDir, data.VideoProgressProperty)
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

	if err := progRdx.ReplaceValues(data.VideoProgressProperty, pr.VideoId, pr.TrimTime()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
