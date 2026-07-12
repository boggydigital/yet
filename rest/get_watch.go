package rest

import (
	"net/http"
	"path"
	"strings"

	"github.com/boggydigital/yet/data"
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

	q := r.URL.Query()

	t := q.Get("t")
	queueDownload := q.Has("queue-download")

	if videoId == "" {
		http.Redirect(w, r, "/list", http.StatusPermanentRedirect)
		return
	}

	// resolve full YouTube URL to just video-id, as needed
	if strings.Contains(videoId, "?") {
		if videoIds, err := yeti.ParseVideoIds(videoId); err != nil {

			// one more attempt - redirect to playlist page if we've got a valid playlist
			if playlistIds, err := yeti.ParsePlaylistIds(videoId); err == nil && len(playlistIds) > 0 {
				http.Redirect(w, r, "/playlist?list="+playlistIds[0], http.StatusPermanentRedirect)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else if len(videoIds) > 0 {
			redirectUrl := path.Join("/watch", videoIds[0])
			if queueDownload {
				redirectUrl += "&queue-download"
			}
			http.Redirect(w, r, redirectUrl, http.StatusPermanentRedirect)
			return
		}
	}

	// iOS insists on inserting a space on paste
	videoId = strings.TrimSpace(videoId)

	var videoIds []string
	if videoIds, err = yeti.ParseVideoIds(videoId); err == nil && len(videoIds) > 0 {
		videoId = videoIds[0]
	}

	if queueDownload {
		if err = rdx.AddValues(data.VideoDownloadQueuedProperty, videoId, yeti.FmtNow()); err != nil {
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

	if err = tmpl.ExecuteTemplate(w, "watch", wvm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
