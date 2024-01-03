package rest

import (
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/rest/view_models"
	"net/http"
)

const (
	showImagesLimit = 20
)

func GetPlaylist(w http.ResponseWriter, r *http.Request) {

	// GET /playlist?id

	var err error
	rdx, err = rdx.RefreshReader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id := r.URL.Query().Get("list")

	if id == "" {
		http.Redirect(w, r, "/new", http.StatusPermanentRedirect)
		return
	}

	pt := view_models.PlaylistTitle(id, rdx)

	w.Header().Set("Content-Type", "text/html")

	// playlist specific styles

	plvm := &view_models.PlaylistViewModel{
		PlaylistId:    id,
		PlaylistTitle: pt,
	}

	if pdq, ok := rdx.GetFirstVal(data.PlaylistDownloadQueueProperty, id); ok && pdq == data.TrueValue {
		plvm.AutoDownloading = true
	}

	if videoIds, ok := rdx.GetAllValues(data.PlaylistVideosProperty, id); ok && len(videoIds) > 0 {
		for i, videoId := range videoIds {
			var options []view_models.VideoOptions
			if i+1 < showImagesLimit {
				options = []view_models.VideoOptions{view_models.ShowPoster, view_models.ShowPublishedDate}
			} else {
				options = []view_models.VideoOptions{view_models.ShowPublishedDate}
			}
			plvm.Videos = append(plvm.Videos, view_models.GetVideoViewModel(videoId, rdx, options...))
		}
	}

	if err := tmpl.ExecuteTemplate(w, "playlist", plvm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
