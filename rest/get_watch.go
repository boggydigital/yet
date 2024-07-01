package rest

import (
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/rest/view_models"
	"github.com/boggydigital/yet/yeti"
	"net/http"
	"strings"
)

func GetWatch(w http.ResponseWriter, r *http.Request) {

	// GET /watch?v&t

	var err error
	rdx, err = rdx.RefreshWriter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	q := r.URL.Query()

	v := q.Get("v")
	t := q.Get("t")
	source := q.Get("source")

	if v == "" {
		http.Redirect(w, r, "/list", http.StatusPermanentRedirect)
		return
	}

	// resolve full YouTube URL to just video-id, as needed
	if strings.Contains(v, "?") {
		if videoIds, err := yeti.ParseVideoIds(v); err != nil {

			// one more attempt - redirect to playlist page if we've got a valid playlist
			if playlistIds, err := yeti.ParsePlaylistIds(v); err == nil && len(playlistIds) > 0 {
				http.Redirect(w, r, "/playlist?list="+playlistIds[0], http.StatusPermanentRedirect)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else if len(videoIds) > 0 {
			redirectUrl := "/watch?v=" + videoIds[0]
			if source != "" {
				redirectUrl += "&source=" + source
			}
			http.Redirect(w, r, redirectUrl, http.StatusPermanentRedirect)
			return
		}
	}

	// iOS insists on inserting a space on paste
	v = strings.TrimSpace(v)

	videoId := ""
	if videoIds, err := yeti.ParseVideoIds(v); err == nil && len(videoIds) > 0 {
		videoId = videoIds[0]
	} else {
		videoId = v
	}

	// set video source unless it's been set already
	if source != "" && !rdx.HasValue(data.VideoSourceProperty, videoId, source) {
		if err := rdx.ReplaceValues(data.VideoSourceProperty, videoId, source); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "text/html")

	wvm, err := view_models.GetWatchViewModel(videoId, t, rdx)
	if err != nil {
		http.Redirect(w, r, "/video_error?v="+videoId+"&err="+err.Error(), http.StatusTemporaryRedirect)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "watch", wvm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
