package rest

import (
	"maps"
	"net/http"
	"path"
	"slices"

	"github.com/boggydigital/redux"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
)

func PostUpdateVideo(w http.ResponseWriter, r *http.Request) {

	// POST /update_video/{videoId}

	var err error
	rdx, err = rdx.RefreshWriter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	videoId := r.PathValue("videoId")

	if err = r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if videoId == "" {
		http.Redirect(w, r, "/list", http.StatusPermanentRedirect)
		return
	}

	boolPropertyInputs := map[string]string{
		data.VideoFavoriteProperty:       "favorite",
		data.VideoForcedDownloadProperty: "forced-download",
	}

	timePropertyInputs := map[string]string{
		data.VideoEndedDateProperty:      "ended",
		data.VideoDownloadQueuedProperty: "download-queued",
	}

	specialProperties := map[string]string{
		data.VideoProgressProperty:    "progress",
		data.VideoEndedReasonProperty: "ended-reason",
	}

	properties := slices.Collect(maps.Keys(boolPropertyInputs))
	properties = append(properties, slices.Collect(maps.Keys(timePropertyInputs))...)
	properties = append(properties, slices.Collect(maps.Keys(specialProperties))...)

	for property, input := range boolPropertyInputs {
		if err = toggleProperty(videoId, property, r.Form.Get(input) == "on", rdx); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	for property, input := range timePropertyInputs {
		if err = toggleTimeProperty(videoId, property, r.Form.Get(input) == "on", rdx); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	for property, input := range specialProperties {
		switch property {
		case data.VideoProgressProperty:
			// progress is cleared (condition: false) when flag IS NOT present in input
			if r.Form.Get(input) != "on" {
				if err = toggleProperty(videoId, property, false, rdx); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			}
		case data.VideoEndedReasonProperty:
			// don't set ended reason unless the video has ended
			if !r.Form.Has("ended") {
				break
			}
			reason := data.DefaultEndedReason
			if er := r.Form.Get(input); er != "" {
				reason = data.ParseVideoEndedReason(er)
			}
			if err = rdx.ReplaceValues(data.VideoEndedReasonProperty, videoId, string(reason)); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
	}

	w.Header().Set("Location", path.Join("/watch", videoId))
	w.WriteHeader(http.StatusSeeOther)
	return
}

func toggleTimeProperty(id, property string, condition bool, rdx redux.Writeable) error {
	if condition {
		return rdx.ReplaceValues(property, id, yeti.FmtNow())
	} else {
		if rdx.HasKey(property, id) {
			return rdx.CutKeys(property, id)
		}
	}
	return nil
}
