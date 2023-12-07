package cli

import (
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yet/yeti"
	"github.com/boggydigital/yt_urls"
	"net/http"
	"net/url"
	"strings"
)

func DownloadHandler(u *url.URL) error {

	ids := strings.Split(u.Query().Get("id"), ",")
	force := u.Query().Has("force")
	return Download(ids, force)
}

func Download(ids []string, force bool) error {

	da := nod.Begin("downloading videos...")
	defer da.End()

	videoIds, err := yeti.ParseVideoIds(ids...)
	if err != nil {
		return da.EndWithError(err)
	}

	metadataDir, err := paths.GetAbsDir(paths.Metadata)
	if err != nil {
		return da.EndWithError(err)
	}

	rxa, err := kvas.ConnectReduxAssets(metadataDir, data.AllProperties()...)
	if err != nil {
		return da.EndWithError(err)
	}

	// adding to the queue before attempting to download
	if err := rxa.BatchAddValues(data.VideosDownloadQueueProperty, trueValues(videoIds...)); err != nil {
		return da.EndWithError(err)
	}

	for _, videoId := range videoIds {

		videoPage, err := yt_urls.GetVideoPage(http.DefaultClient, videoId)
		if err != nil {
			return da.EndWithError(err)
		}

		relFilename := yeti.DefaultFilenameDelegate(videoId, videoPage)

		if err := downloadVideo(dolo.DefaultClient, relFilename, videoPage); err != nil {
			return da.EndWithError(err)
		}

		if err := getVideoPageMetadata(videoPage, videoId, rxa); err != nil {
			return da.EndWithError(err)
		}

		if err := yeti.GetPosters(videoPage, dolo.DefaultClient); err != nil {
			return da.EndWithError(err)
		}

		if err := getVideoPageCaptions(videoPage, videoId, rxa, dolo.DefaultClient); err != nil {
			return da.EndWithError(err)
		}

		// remove from the queue upon successful download
		if err := rxa.CutVal(data.VideosDownloadQueueProperty, videoId, data.TrueValue); err != nil {
			return da.EndWithError(err)
		}

		// add to watchlist upon successful download
		if err := rxa.AddValues(data.VideosWatchlistProperty, videoId, data.TrueValue); err != nil {
			return da.EndWithError(err)
		}

	}

	da.EndWithResult("done")

	return nil
}
