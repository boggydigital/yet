package rest

import (
	"github.com/boggydigital/yet/data"
	"io"
	"net/http"
	"strings"
)

const (
	showImagesLimit = 20
)

func GetPlaylist(w http.ResponseWriter, r *http.Request) {

	// GET /playlist?id

	var err error
	rdx, err = rdx.RefreshWriter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id := r.URL.Query().Get("list")

	if id == "" {
		http.Redirect(w, r, "/new", http.StatusPermanentRedirect)
		return
	}

	pt := playlistTitle(id, rdx)

	w.Header().Set("Content-Type", "text/html")

	sb := &strings.Builder{}
	sb.WriteString("<!doctype html>")
	sb.WriteString("<html>")
	sb.WriteString("<head>" +
		"<meta charset='UTF-8'>" +
		"<link rel='icon' href='data:image/svg+xml,<svg xmlns=%22http://www.w3.org/2000/svg%22 viewBox=%220 0 100 100%22><text y=%22.9em%22 font-size=%2290%22>🔻</text></svg>' type='image/svg+xml'/>" +
		"<meta name='viewport' content='width=device-width, initial-scale=1.0'>" +
		"<meta name='color-scheme' content='dark light'>" +
		"<title>🔻 " + pt + "</title>" +
		"<style>")

	writeSharedStyles(sb)

	// playlist specific styles
	sb.WriteString("a.refresh {color: dodgerblue; margin-block: 2rem;}")

	sb.WriteString("</style></head>")
	sb.WriteString("<body>")

	sb.WriteString("<h1><span class='playlistTitle'>" + pt + "</span></h1>")

	if pdq, ok := rdx.GetFirstVal(data.PlaylistDownloadQueueProperty, id); ok && pdq == data.TrueValue {
		sb.WriteString("<div class='subtle'>Automatically refreshing and downloading new videos</div>")
	}

	sb.WriteString("<a class='video refresh' href='/refresh?list=" + id + "'>Refresh playlist</a>")

	if videoIds, ok := rdx.GetAllValues(data.PlaylistVideosProperty, id); ok && len(videoIds) > 0 {
		for i, videoId := range videoIds {
			var options []VideoOptions
			if i+1 < showImagesLimit {
				options = []VideoOptions{ShowPoster, ShowPublishedDate}
			} else {
				options = []VideoOptions{ShowPublishedDate}
			}
			writeVideo(videoId, rdx, sb, options...)
		}
	}

	sb.WriteString("</body>")
	sb.WriteString("</html>")

	if _, err := io.WriteString(w, sb.String()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
