package rest

import (
	"bytes"
	_ "embed"
	"github.com/boggydigital/yet/paths"
	"io"
	"net/http"
	"os"
	"time"
)

//go:embed "posters/yet_maxresdefault.png"
var yetPosterMaxResDefault []byte

//go:embed "posters/yet_hqdefault.png"
var yetPosterHQDefault []byte

func GetPoster(w http.ResponseWriter, r *http.Request) {

	// GET /poster?v&q

	videoId := r.URL.Query().Get("v")
	quality := r.URL.Query().Get("q")

	absPosterFilename, err := paths.AbsPosterPath(videoId, quality)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := os.Stat(absPosterFilename); err == nil {
		http.ServeFile(w, r, absPosterFilename)
	} else {

		var br io.ReadSeeker
		filename := ""

		switch quality {
		case "maxresdefault":
			filename = "yet_maxresdefault.png"
			br = bytes.NewReader(yetPosterMaxResDefault)
		case "hqdefault":
			filename = "yet_hqdefault.png"
			br = bytes.NewReader(yetPosterHQDefault)
		}

		if br != nil {
			http.ServeContent(w, r, filename, time.Unix(0, 0), br)
		} else {
			http.NotFound(w, r)
		}

		return
	}

}
