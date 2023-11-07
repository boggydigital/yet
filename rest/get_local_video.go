package rest

import (
	"github.com/boggydigital/yet/paths"
	"net/http"
	"os"
	"path/filepath"
)

func GetLocalVideo(w http.ResponseWriter, r *http.Request) {

	// GET /local_video?file

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
		http.Error(w, "local_video requires local file", http.StatusBadRequest)
		return
	}
}
