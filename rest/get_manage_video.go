package rest

import (
	"net/http"
)

func GetManageVideo(w http.ResponseWriter, r *http.Request) {

	// GET /manage_video/{videoId}

	var err error
	rdx, err = rdx.RefreshWriter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	videoId := r.PathValue("videoId")

	if videoId == "" {
		http.Redirect(w, r, "/list", http.StatusPermanentRedirect)
		return
	}

	//w.Header().Set("Content-Type", "text/html")
	//
	//if err := tmpl.ExecuteTemplate(w, "manage_video", view_models.GetVideoManagementModel(videoId, rdx)); err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}

}
