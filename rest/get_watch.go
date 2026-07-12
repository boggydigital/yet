package rest

import (
	"net/http"
	"strings"

	"github.com/boggydigital/yet/rest/view_models"
	"github.com/boggydigital/yet/yeti"
)

func GetWatch(w http.ResponseWriter, r *http.Request) {

	// GET /watch/{videoId}?t

	var err error
	rdx, err = rdx.RefreshWriter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	videoId := r.PathValue("videoId")

	t := r.URL.Query().Get("t")

	if videoId == "" {
		http.Redirect(w, r, "/list", http.StatusPermanentRedirect)
		return
	}

	// iOS insists on inserting a space on paste
	videoId = strings.TrimSpace(videoId)

	var videoIds []string
	if videoIds, err = yeti.ParseVideoIds(videoId); err == nil && len(videoIds) > 0 {
		videoId = videoIds[0]
	}

	w.Header().Set("Content-Type", "text/html")

	wvm, err := view_models.GetWatchViewModel(videoId, t, rdx)
	if err != nil {
		http.Redirect(w, r, "/video_error?v="+videoId+"&err="+err.Error(), http.StatusTemporaryRedirect)
		return
	}

	if err = tmpl.ExecuteTemplate(w, "watch", wvm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
