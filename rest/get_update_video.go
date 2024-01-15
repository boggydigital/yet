package rest

import (
	"fmt"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/pasu"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"net/http"
	"net/url"
	"time"
)

func GetUpdateVideo(w http.ResponseWriter, r *http.Request) {

	// GET /update_video?v

	videoId := r.URL.Query().Get("v")

	if videoId == "" {
		http.Redirect(w, r, "/list", http.StatusPermanentRedirect)
		return
	}

	metadataDir, err := pasu.GetAbsDir(paths.Metadata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	properties := []string{data.VideoProgressProperty,
		data.VideoEndedProperty,
		data.VideosWatchlistProperty,
		data.VideosDownloadQueueProperty}

	vRdx, err := kvas.NewReduxWriter(metadataDir, properties...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, p := range properties {
		if err := updateVideoProperty(videoId, p, r.URL, vRdx); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	http.Redirect(w, r, "/watch?v="+videoId, http.StatusTemporaryRedirect)
}

func updateVideoProperty(videoId string, property string, u *url.URL, rdx kvas.WriteableRedux) error {

	flagStr := ""
	switch property {
	case data.VideoProgressProperty:
		flagStr = "progress"
	case data.VideoEndedProperty:
		flagStr = "ended"
	case data.VideosWatchlistProperty:
		flagStr = "watchlist"
	case data.VideosDownloadQueueProperty:
		flagStr = "download"
	default:
		return fmt.Errorf("unsupported property %s", property)
	}

	flag := u.Query().Has(flagStr)

	var err error

	if flag {
		if !rdx.HasKey(property, videoId) {
			value := data.TrueValue
			switch property {
			case data.VideoProgressProperty:
				// setting progress requires current time - users should be encouraged to scrub video instead
				return nil
			case data.VideoEndedProperty:
				// ended requires current time as a value to set
				value = time.Now().Format(time.RFC3339)
			}
			err = rdx.AddValues(property, videoId, value)
		}
	} else {
		if rdx.HasKey(property, videoId) {
			err = rdx.CutKeys(property, videoId)
		}
	}

	return err
}
