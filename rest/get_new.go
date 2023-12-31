package rest

import (
	"net/http"
)

func GetNew(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")

	if err := tmpl.ExecuteTemplate(w, "new", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
