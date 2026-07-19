package rest

import (
	"net/http"
	"path"

	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
)

func GetEnded(w http.ResponseWriter, r *http.Request) {

	// GET /ended/{video}/{reason}

	var err error
	rdx, err = rdx.RefreshWriter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	videoId := r.PathValue("videoId")
	reason := data.ParseVideoEndedReason(r.PathValue("reason"))

	if videoId != "" {

		// store completion timestamp
		if err = rdx.ReplaceValues(data.VideoEndedDateProperty, videoId, yeti.FmtNow()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// store ended reason if not-default
		if reason != data.DefaultEndedReason {
			if err = rdx.ReplaceValues(data.VideoEndedReasonProperty, videoId, reason.String()); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	http.Redirect(w, r, path.Join("/watch", videoId), http.StatusTemporaryRedirect)
}
