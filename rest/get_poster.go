package rest

import (
	"bytes"
	_ "embed"
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yet/yeti"
	"github.com/boggydigital/yt_urls"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

//go:embed "posters/yet_maxresdefault.png"
var yetPosterMaxResDefault []byte

//go:embed "posters/yet_hqdefault.png"
var yetPosterHQDefault []byte

func GetPoster(w http.ResponseWriter, r *http.Request) {

	// GET /poster?v&q

	videoId := r.URL.Query().Get("v")
	tq := r.URL.Query().Get("q")

	quality := yt_urls.ParseThumbnailQuality(tq)

	if videoId == "" {
		return
	}

	for q := quality; q != yt_urls.ThumbnailQualityUnknown; q = yt_urls.LowerQuality(q) {

		absPosterFilename, err := paths.AbsPosterPath(videoId, q)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// attempt to fetch posters from the origin if they don't exist locally
		// unless that's a URL file
		if _, err := os.Stat(absPosterFilename); os.IsNotExist(err) &&
			!strings.Contains(videoId, yt_urls.DefaultVideoExt) {
			if err := yeti.GetPosters(videoId, dolo.DefaultClient, q); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		if _, err := os.Stat(absPosterFilename); err == nil {
			http.ServeFile(w, r, absPosterFilename)
			return
		}
	}

	var br io.ReadSeeker
	filename := ""

	switch quality {
	case yt_urls.ThumbnailQualityMaxRes:
		filename = "yet_maxresdefault.png"
		br = bytes.NewReader(yetPosterMaxResDefault)
	default:
		filename = "yet_hqdefault.png"
		br = bytes.NewReader(yetPosterHQDefault)
	}

	if br != nil {
		http.ServeContent(w, r, filename, time.Unix(0, 0), br)
	} else {
		http.NotFound(w, r)
	}

}
