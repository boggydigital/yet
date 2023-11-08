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
	sb.WriteString("<head><style>" +
		"body {background: black; color: white;font-family:sans-serif; margin: 1rem;} " +
		"a {display:block;color:lightblue;font-size:1.3rem;font-weight:bold;text-decoration:none;margin-block:0.5rem}" +
		"</style></head>")
	sb.WriteString("<body>")

	for _, videoId := range rxa.Keys(data.VideoTitleProperty) {

		if title, ok := rxa.GetFirstVal(data.VideoTitleProperty, videoId); ok {
			sb.WriteString("<a href='/watch?v=" + videoId + "'>" + title + "</a>")
		}
	}

	sb.WriteString("</body>")
	sb.WriteString("</html>")

	if _, err := io.WriteString(w, sb.String()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
