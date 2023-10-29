package rest

import (
	"io"
	"net/http"
)

func GetWatch(w http.ResponseWriter, r *http.Request) {

	// GET /watch?video-id

	if _, err := io.WriteString(w, "watch"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
