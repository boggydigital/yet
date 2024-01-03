package rest

import "net/http"

func GetSearch(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")

	if err := tmpl.ExecuteTemplate(w, "search", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
