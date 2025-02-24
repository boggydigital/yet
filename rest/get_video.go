package rest

import (
	"github.com/boggydigital/pathways"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet_urls/youtube_urls"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func GetVideo(w http.ResponseWriter, r *http.Request) {

	// GET /video?file

	file := r.URL.Query().Get("file")

	if filepath.IsLocal(file) {

		absVideosDir, err := pathways.GetAbsDir(data.Videos)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		absFilepath := filepath.Join(absVideosDir, file)

		if _, err := os.Stat(absFilepath); err == nil {
			if strings.HasSuffix(file, youtube_urls.DefaultVideoExt) {
				w.Header().Set("Content-Type", "video/mp4")
			}
			http.ServeFile(w, r, absFilepath)
		} else {
			http.NotFound(w, r)
		}
	} else {
		http.Error(w, "file is not local to server videos dir", http.StatusBadRequest)
		return
	}
}
