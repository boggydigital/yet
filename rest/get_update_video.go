package rest

import (
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
	"golang.org/x/exp/maps"
	"net/http"
)

func GetUpdateVideo(w http.ResponseWriter, r *http.Request) {

	// GET /update_video?v

	var err error
	rdx, err = rdx.RefreshWriter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	q := r.URL.Query()

	videoId := q.Get("v")

	if videoId == "" {
		http.Redirect(w, r, "/list", http.StatusPermanentRedirect)
		return
	}

	boolPropertyInputs := map[string]string{
		data.VideoFavoriteProperty:           "favorite",
		data.VideoForcedDownloadProperty:     "forced-download",
		data.VideoPreferSingleFormatProperty: "prefer-single-format",
	}

	timePropertyInputs := map[string]string{
		data.VideoEndedDateProperty:      "ended",
		data.VideoDownloadQueuedProperty: "download-queued",
	}

	specialProperties := map[string]string{
		data.VideoProgressProperty:    "progress",
		data.VideoEndedReasonProperty: "ended-reason",
		data.VideoSourceProperty:      "clear-source",
	}

	properties := maps.Keys(boolPropertyInputs)
	properties = append(properties, maps.Keys(timePropertyInputs)...)
	properties = append(properties, maps.Keys(specialProperties)...)

	for property, input := range boolPropertyInputs {
		if err := toggleProperty(videoId, property, q.Has(input), rdx); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	for property, input := range timePropertyInputs {
		if err := toggleTimeProperty(videoId, property, q.Has(input), rdx); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	for property, input := range specialProperties {
		switch property {
		case data.VideoProgressProperty:
			// progress and source and handled in the same way:
			// - nothing happens on set
			// - the value is remove on clear
			fallthrough
		case data.VideoSourceProperty:
			if q.Has(input) {
				if err := toggleProperty(videoId, property, false, rdx); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			} else {
				// do nothing, clear source if only for clearing sources
			}
		case data.VideoEndedReasonProperty:
			// don't set ended reason unless the video has ended
			if !q.Has("ended") {
				break
			}
			reason := data.DefaultEndedReason
			if er := q.Get(input); er != "" {
				reason = data.ParseVideoEndedReason(er)
			}
			if err := rdx.ReplaceValues(data.VideoEndedReasonProperty, videoId, string(reason)); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
	}

	http.Redirect(w, r, "/watch?v="+videoId, http.StatusTemporaryRedirect)
}

func toggleTimeProperty(id, property string, condition bool, rdx kevlar.WriteableRedux) error {
	if condition {
		return rdx.ReplaceValues(property, id, yeti.FmtNow())
	} else {
		if rdx.HasKey(property, id) {
			return rdx.CutKeys(property, id)
		}
	}
	return nil
}
