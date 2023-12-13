package rest

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"io"
	"net/http"
	"strings"
)

func GetPlaylist(w http.ResponseWriter, r *http.Request) {

	// GET /playlist?id

	id := r.URL.Query().Get("id")

	if id == "" {
		http.Redirect(w, r, "/new", http.StatusPermanentRedirect)
		return
	}

	absMetadataDir, err := paths.GetAbsDir(paths.Metadata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rxa, err := kvas.ConnectReduxAssets(absMetadataDir, data.AllProperties()...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	sb := &strings.Builder{}
	sb.WriteString("<!doctype html>")
	sb.WriteString("<html>")
	sb.WriteString("<head>" +
		"<meta charset='UTF-8'>" +
		"<link rel='icon' href='data:image/svg+xml,<svg xmlns=%22http://www.w3.org/2000/svg%22 viewBox=%220 0 100 100%22><text y=%22.9em%22 font-size=%2290%22>🔻</text></svg>' type='image/svg+xml'/>" +
		"<meta name='viewport' content='width=device-width, initial-scale=1.0'>" +
		"<meta name='color-scheme' content='dark light'>" +
		"<style>" +
		"body {background: black; color: white;font-family:sans-serif; margin: 1rem;} " +
		"a.video {display:block;color:white;font-size:1.3rem;font-weight:bold;text-decoration:none;margin-block:0.5rem;margin-block-end: 1rem}" +
		"a.video img {border-radius:0.25rem;width:200px}" +
		"a.video span {font-size:1rem}" +
		"</style></head>")
	sb.WriteString("<body>")

	if plt, ok := rxa.GetFirstVal(data.PlaylistTitleProperty, id); ok && plt != "" {
		sb.WriteString("<h1>" + plt + "</h1>")
	}

	if videoIds, ok := rxa.GetAllValues(data.PlaylistVideosProperty, id); ok && len(videoIds) > 0 {
		for _, videoId := range videoIds {
			writeVideo(videoId, rxa, sb)
		}
	}

	sb.WriteString("</body>")
	sb.WriteString("</html>")

	if _, err := io.WriteString(w, sb.String()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}