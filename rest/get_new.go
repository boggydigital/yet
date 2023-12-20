package rest

import (
	"io"
	"net/http"
	"strings"
)

func GetNew(w http.ResponseWriter, r *http.Request) {

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
		"input[type='text'] {width:90%}" +
		"input {font-size:1.25rem;display:block;}" +
		"input[type='submit'] {margin-block: 1rem;}" +
		"a.video {display:block;color:white;font-size:1.3rem;font-weight:bold;text-decoration:none;margin-block:0.5rem;margin-block-end: 1rem}" +
		"a.highlight {color:gold; margin-block:2rem}" +
		"</style></head>")
	sb.WriteString("<body>")

	sb.WriteString("<a class='video highlight' href='/list'>Watch list</a>")

	sb.WriteString("<form method='get' action='/watch'>")
	sb.WriteString("<input id='v' name='v' type='text' placeholder='Paste or enter YouTube link or video-id' />")
	sb.WriteString("<input type='submit' />")
	sb.WriteString("</form>")

	sb.WriteString("</body>")
	sb.WriteString("</html>")

	if _, err := io.WriteString(w, sb.String()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
