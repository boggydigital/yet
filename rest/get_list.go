package rest

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"io"
	"net/http"
	"strings"
)

func GetList(w http.ResponseWriter, r *http.Request) {

	// GET /list

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
		"a {display:block;color:lightblue;font-size:1.3rem;font-weight:bold;text-decoration:none;margin-block:0.5rem}" +
		"</style></head>")
	sb.WriteString("<body>")

	// continue watching
	// videos watchlist
	// videos download queue

	sb.WriteString("<h1>Continue watching</h1>")
	for _, id := range rxa.Keys(data.VideoProgressProperty) {
		if ended, ok := rxa.GetFirstVal(data.VideoEndedProperty, id); !ok || ended == "" {
			writeVideo(id, rxa, sb)
		}
	}

	sb.WriteString("<h1>Watchlist</h1>")
	for _, id := range rxa.Keys(data.VideosWatchlistProperty) {
		if le, ok := rxa.GetFirstVal(data.VideoEndedProperty, id); ok && le != "" {
			continue
		}
		writeVideo(id, rxa, sb)
	}

	sb.WriteString("<h1>URL Watchlist</h1>")
	for _, id := range rxa.Keys(data.UrlsWatchlistProperty) {
		writeVideo(id, rxa, sb)
	}

	sb.WriteString("<h1>Download queue</h1>")
	for _, id := range rxa.Keys(data.VideosDownloadQueueProperty) {
		writeVideo(id, rxa, sb)
	}

	sb.WriteString("</body>")
	sb.WriteString("</html>")

	if _, err := io.WriteString(w, sb.String()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func writeVideo(videoId string, rxa kvas.ReduxAssets, sb *strings.Builder) {

	videoTitle := videoId
	if title, ok := rxa.GetFirstVal(data.VideoTitleProperty, videoId); ok && title != "" {
		videoTitle = title
	}

	sb.WriteString("<a href='/watch?v=" + videoId + "'>" + videoTitle + "</a>")

}
