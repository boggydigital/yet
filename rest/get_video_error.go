package rest

import (
	"github.com/boggydigital/yet/rest/view_models"
	"net/http"
)

func GetVideoError(w http.ResponseWriter, r *http.Request) {

	// GET /video_error?v&err

	videoId := r.URL.Query().Get("v")
	errStr := r.URL.Query().Get("err")

	if err := tmpl.ExecuteTemplate(w, "video_error", view_models.GetVideoErrorViewModel(videoId, errStr, rdx)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
