package rest

import (
	"encoding/json"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/pasu"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"net/http"
	"time"
)

type EndedRequest struct {
	VideoId string `json:"v"`
}

func PostEnded(w http.ResponseWriter, r *http.Request) {

	// POST /ended
	// {v}

	metadataDir, err := pasu.GetAbsDir(paths.Metadata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	endRdx, err := kvas.NewReduxWriter(metadataDir, data.VideoEndedProperty, data.PlaylistNewVideosProperty)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	currentTime := time.Now().Format(time.RFC3339)
	if err := endRdx.ReplaceValues(data.VideoEndedProperty, er.VideoId, currentTime); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// remove the video from playlist new videos
	for _, playlistId := range endRdx.Keys(data.PlaylistNewVideosProperty) {
		if endRdx.HasValue(data.PlaylistNewVideosProperty, playlistId, er.VideoId) {
			if err := endRdx.CutValues(data.PlaylistNewVideosProperty, playlistId, er.VideoId); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
}
