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

	// store completion timestamp
	currentTime := time.Now().Format(http.TimeFormat)
	if err := rdx.ReplaceValues(data.VideoEndedProperty, er.VideoId, currentTime); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// remove the video from playlist new videos
	for _, playlistId := range rdx.Keys(data.PlaylistNewVideosProperty) {
		if rdx.HasValue(data.PlaylistNewVideosProperty, playlistId, er.VideoId) {
			if err := rdx.CutValues(data.PlaylistNewVideosProperty, playlistId, er.VideoId); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
}
