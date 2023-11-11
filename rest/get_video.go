package rest

import (
	"github.com/boggydigital/yet/paths"
	"net/http"
	"os"
	"path/filepath"
)

func GetVideo(w http.ResponseWriter, r *http.Request) {

	// GET /video?file

	file := r.URL.Query().Get("file")

	if filepath.IsLocal(file) {

		absVideosDir, err := paths.GetAbsDir(paths.Videos)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		absFilepath := filepath.Join(absVideosDir, file)

		if _, err := os.Stat(absFilepath); err == nil {
			http.ServeFile(w, r, absFilepath)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "file is not local to server videos dir", http.StatusBadRequest)
		return
	}
}
