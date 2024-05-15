package rest

import (
	"github.com/boggydigital/pathways"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yt_urls"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func GetVideo(w http.ResponseWriter, r *http.Request) {

	// GET /video?file

	file := r.URL.Query().Get("file")

	if filepath.IsLocal(file) {

		absVideosDir, err := pathways.GetAbsDir(paths.Videos)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		absFilepath := filepath.Join(absVideosDir, file)

		if _, err := os.Stat(absFilepath); err == nil {
			if strings.HasSuffix(file, yt_urls.DefaultVideoExt) {
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
