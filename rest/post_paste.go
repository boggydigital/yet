package rest

import (
	"net/http"
	"path"

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

	videoId := r.FormValue(paramVideoId)

	// resolve full YouTube URL to just video-id, as needed
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

	queueDownload := r.FormValue(paramQueueDownload) == "on"
	downloadVideo := r.FormValue(paramDownloadVideo) == "on"

	if downloadVideo {
		http.Redirect(w, r, path.Join("/download_video", videoId), http.StatusTemporaryRedirect)
		return
	}

	if queueDownload {
		http.Redirect(w, r, path.Join("/queue_download", videoId), http.StatusTemporaryRedirect)
		return
	}

	// when neither download nor queue download are requested - redirect to watch page

	w.Header().Set("Location", path.Join("/watch", videoId))
	w.WriteHeader(http.StatusSeeOther)

	return
}
