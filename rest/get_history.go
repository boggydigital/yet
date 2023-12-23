package rest

import (
	"github.com/boggydigital/yet/data"
	"io"
	"net/http"
	"strings"
)

func GetHistory(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")

	sb := &strings.Builder{}
	sb.WriteString("<!doctype html>")
	sb.WriteString("<html>")
	sb.WriteString("<head>" +
		"<meta charset='UTF-8'>" +
		"<link rel='icon' href='data:image/svg+xml,<svg xmlns=%22http://www.w3.org/2000/svg%22 viewBox=%220 0 100 100%22><text y=%22.9em%22 font-size=%2290%22>🔻</text></svg>' type='image/svg+xml'/>" +
		"<meta name='viewport' content='width=device-width, initial-scale=1.0'>" +
		"<meta name='color-scheme' content='dark light'>" +
		"<title>🔻 History</title>" +
		"<style>" +
		"body {background: black; color: white;font-family:sans-serif; margin: 1rem;} " +
		"a.video {display:flex;flex-direction: column; color:white;font-size:1.3rem;font-weight:bold;text-decoration:none;margin-block:0.5rem;margin-block-end: 1rem}" +
		"a.video img {border-radius:0.25rem;width:200px}" +
		"a.video.ended {filter:grayscale(1.0)}" +
		"a.highlight {color:gold; margin-block:2rem}" +
		"details {margin-block:2rem; content-visibility: auto}" +
		"summary {margin-block-end: 2rem}" +
		"summary h1 {display: inline; cursor: pointer; margin-inline-start: 0.5rem;color:turquoise}" +
		"a.playlist {display:block;color:deeppink;font-size:1.3rem;font-weight:bold;text-decoration:none;margin-block:0.5rem;margin-block-end: 1rem}" +
		"ul {list-style:none; padding-inline-start: 1rem}" +
		".title {margin-block-start:0.5rem;margin-block-end:0.25rem}" +
		".subtitle {font-size:62.5%; font-weight:normal}" +
		"a.video .subtitle {color:dimgray}" +
		"</style></head>")
	sb.WriteString("<body>")

	sb.WriteString("<h1>Watch history</h1>")

	var err error

	whKeys := rdx.Keys(data.VideoEndedProperty)
	if len(whKeys) > 0 {

		whKeys, err = rdx.Sort(whKeys, false, data.VideoTitleProperty)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, id := range whKeys {
			writeVideo(id, false, rdx, sb)
		}
	}

	sb.WriteString("</body>")
	sb.WriteString("</html>")

	if _, err := io.WriteString(w, sb.String()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
