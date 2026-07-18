package rest

import (
	"net/http"
	"path"
	"strings"

	"github.com/boggydigital/yet/yeti"
)

func GetPasteVideo(w http.ResponseWriter, r *http.Request) {

	// GET /paste_video/?videoId&download-video&queue-download

	var err error
	rdx, err = rdx.RefreshWriter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	q := r.URL.Query()

	videoId := q.Get("video-id")

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

	downloadVideo := q.Has("download-video")
	queueDownload := q.Has("queue-download")

	if downloadVideo {
		http.Redirect(w, r, path.Join("/download_video", videoId), http.StatusTemporaryRedirect)
		return
	}

	if queueDownload {
		http.Redirect(w, r, path.Join("/queue_download", videoId), http.StatusTemporaryRedirect)
		return
	}

	//if !downloadVideo && !queueDownload {
	http.Redirect(w, r, "/video_error?v="+videoId+"&err=Paste+requires+Download+now+or+Queue+download", http.StatusTemporaryRedirect)
	return
}
