package rest

import (
	"encoding/json"
	"github.com/boggydigital/kvas"
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

	var err error
	progressRdx, err = progressRdx.RefreshWriter()
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
	currentTime := time.Now().Format(time.RFC3339)
	if err := progressRdx.ReplaceValues(data.VideoEndedProperty, er.VideoId, currentTime); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// remove the video from playlist new videos
	if err := rmVideoFromPlaylistNewVideos(er.VideoId, progressRdx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func rmVideoFromPlaylistNewVideos(videoId string, rdx kvas.WriteableRedux) error {
	if err := rdx.MustHave(data.PlaylistNewVideosProperty); err != nil {
		return err
	}
	for _, playlistId := range progressRdx.Keys(data.PlaylistNewVideosProperty) {
		if progressRdx.HasValue(data.PlaylistNewVideosProperty, playlistId, videoId) {
			if err := progressRdx.CutValues(data.PlaylistNewVideosProperty, playlistId, videoId); err != nil {
				return err
			}
		}
	}
	return nil
}
