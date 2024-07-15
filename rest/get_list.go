package rest

import (
	"github.com/boggydigital/yet/rest/view_models"
	"net/http"
)

func GetList(w http.ResponseWriter, r *http.Request) {

	// GET /list

	var err error
	rdx, err = rdx.RefreshWriter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	lvm, err := view_models.GetListViewModel(rdx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "list", lvm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
