package rest

import (
	"fmt"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/pathways"
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

	metadataDir, err := pathways.GetAbsDir(paths.Metadata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	properties := []string{
		data.VideoProgressProperty,
		data.VideoEndedDateProperty,
		data.VideoEndedReasonProperty,
		//data.VideosWatchlistProperty,
		data.VideoDownloadQueuedProperty,
		data.VideoForcedDownloadProperty,
		data.VideoPreferSingleFormatProperty,
		//data.PlaylistNewVideosProperty,
	}

	vRdx, err := kvas.NewReduxWriter(metadataDir, properties...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, p := range properties {
		//if p == data.PlaylistNewVideosProperty {
		//	continue
		//}
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
	case data.VideoEndedDateProperty:
		flagStr = "ended"
	case data.VideoEndedReasonProperty:
		flagStr = "ended-reason"
	//case data.VideosWatchlistProperty:
	//	flagStr = "watchlist"
	case data.VideoDownloadQueuedProperty:
		flagStr = "download"
	case data.VideoForcedDownloadProperty:
		flagStr = "forced-download"
	case data.VideoPreferSingleFormatProperty:
		flagStr = "single-format"
	default:
		return fmt.Errorf("unsupported property %s", property)
	}

	flag := u.Query().Has(flagStr)

	if flag {
		if !rdx.HasKey(property, videoId) {
			value := data.TrueValue
			switch property {
			case data.VideoProgressProperty:
				// setting progress requires current time - users should be encouraged to scrub video instead
				return nil
			case data.VideoEndedReasonProperty:
				// skipped video must also make sure the video is set as ended
				if !rdx.HasKey(data.VideoEndedDateProperty, videoId) {
					t := time.Now().Format(time.RFC3339)
					if err := rdx.AddValues(data.VideoEndedDateProperty, videoId, t); err != nil {
						return err
					}
				}
			case data.VideoEndedDateProperty:
				// ended requires current time as a value to set
				value = time.Now().Format(time.RFC3339)
			}
			if err := rdx.AddValues(property, videoId, value); err != nil {
				return err
			}

			if property == data.VideoEndedReasonProperty || property == data.VideoEndedDateProperty {
				// if ended or skipped - remove video from new playlist videos
				if err := rmVideoFromPlaylistNewVideos(videoId, rdx); err != nil {
					return err
				}
			}
		}
	} else {
		if rdx.HasKey(property, videoId) {
			if err := rdx.CutKeys(property, videoId); err != nil {
				return err
			}
		}
	}

	return nil
}
