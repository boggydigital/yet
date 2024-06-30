package rest

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/pathways"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yet/yeti"
	"golang.org/x/exp/maps"
	"net/http"
)

func GetUpdateVideo(w http.ResponseWriter, r *http.Request) {

	// GET /update_video?v

	q := r.URL.Query()

	videoId := q.Get("video-id")

	if videoId == "" {
		http.Redirect(w, r, "/list", http.StatusPermanentRedirect)
		return
	}

	metadataDir, err := pathways.GetAbsDir(paths.Metadata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	boolPropertyInputs := map[string]string{
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
	}

	properties := maps.Keys(boolPropertyInputs)
	properties = append(properties, maps.Keys(timePropertyInputs)...)
	properties = append(properties, maps.Keys(specialProperties)...)

	vRdx, err := kvas.NewReduxWriter(metadataDir, properties...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for property, input := range boolPropertyInputs {
		if err := toggleProperty(videoId, property, q.Has(input), vRdx); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	for property, input := range timePropertyInputs {
		if err := toggleTimeProperty(videoId, property, q.Has(input), vRdx); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	for property, input := range specialProperties {
		switch property {
		case data.VideoProgressProperty:
			if q.Has(input) {
				// do nothing, progress cannot be set
			} else {
				if err := toggleProperty(videoId, property, q.Has(input), vRdx); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			}
		case data.VideoEndedReasonProperty:
			reason := data.DefaultEndedReason
			if er := q.Get(input); er != "" {
				reason = data.ParseVideoEndedReason(er)
			}
			if err := vRdx.ReplaceValues(data.VideoEndedReasonProperty, videoId, string(reason)); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

	}

	http.Redirect(w, r, "/watch?v="+videoId, http.StatusTemporaryRedirect)
}

func toggleTimeProperty(id, property string, condition bool, rdx kvas.WriteableRedux) error {
	if condition {
		return rdx.ReplaceValues(property, id, yeti.FmtNow())
	} else {
		if rdx.HasKey(property, id) {
			return rdx.CutKeys(property, id)
		}
	}
	return nil
}
