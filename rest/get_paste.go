package rest

import (
	"net/http"
)

func GetPaste(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")

	if err := tmpl.ExecuteTemplate(w, "paste", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
