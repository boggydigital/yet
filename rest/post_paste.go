package rest

import (
	"net/http"
	"path"
	"strings"

	"github.com/boggydigital/yet/yeti"
)

func PostPaste(w http.ResponseWriter, r *http.Request) {

	// POST /paste

	var err error
	rdx, err = rdx.RefreshWriter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	videoId := r.FormValue("video-id")

	// resolve full YouTube URL to just video-id, as needed
	if strings.Contains(videoId, "?") {
		var videoIds []string
		if videoIds, err = yeti.ParseVideoIds(videoId); err != nil {

			// one more attempt - redirect to playlist page if we've got a valid playlist
			var playlistIds []string
			if playlistIds, err = yeti.ParsePlaylistIds(videoId); err == nil && len(playlistIds) > 0 {
				http.Redirect(w, r, path.Join("/playlist", playlistIds[0]), http.StatusPermanentRedirect)
				return
			}

			return
		} else if len(videoIds) > 0 {
			videoId = videoIds[0]
		}
	}

	downloadVideo := r.FormValue("download-video") == "on"
	queueDownload := r.FormValue("queue-download") == "on"

	if downloadVideo {
		w.Header().Set("Location", path.Join("/download_video", videoId))
		w.WriteHeader(http.StatusSeeOther)
		return
	}

	if queueDownload {
		w.Header().Set("Location", path.Join("/queue_download", videoId))
		w.WriteHeader(http.StatusSeeOther)
		return
	}

	// when neither download nor queue download are requested - redirect to watch page

	w.Header().Set("Location", path.Join("/watch", videoId))
	w.WriteHeader(http.StatusSeeOther)

	return
}
