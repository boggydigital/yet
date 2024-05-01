package cli

import (
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pasu"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yet/yeti"
	"github.com/boggydigital/yt_urls"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func DownloadHandler(u *url.URL) error {

	ids := strings.Split(u.Query().Get("id"), ",")
	queue := u.Query().Has("queue")
	force := u.Query().Has("force")
	singleFormat := u.Query().Has("single-format")
	return Download(ids, queue, force, singleFormat)
}

func Download(ids []string, queue, force, singleFormat bool) error {

	da := nod.NewProgress("downloading videos...")
	defer da.End()

	metadataDir, err := pasu.GetAbsDir(paths.Metadata)
	if err != nil {
		return da.EndWithError(err)
	}

	rdx, err := kvas.NewReduxWriter(metadataDir, data.AllProperties()...)
	if err != nil {
		return da.EndWithError(err)
	}

	if queue {
		ids = append(ids, rdx.Keys(data.VideosDownloadQueueProperty)...)
	}

	videoIds, err := yeti.ParseVideoIds(ids...)
	if err != nil {
		return da.EndWithError(err)
	}

	da.TotalInt(len(videoIds))

	// adding to the queue before attempting to download
	if err := rdx.BatchAddValues(data.VideosDownloadQueueProperty, trueValues(videoIds...)); err != nil {
		return da.EndWithError(err)
	}

	for _, videoId := range videoIds {

		videoForce := rdx.HasKey(data.VideoForcedDownloadProperty, videoId) || force
		videoSingleFormat := rdx.HasKey(data.VideoSingleFormatDownloadProperty, videoId) || singleFormat

		videoPage, err := yeti.GetVideoPage(videoId)
		if err != nil {
			da.Error(err)
			continue
		}

		if err := yeti.DecodeSignatureCiphers(http.DefaultClient, videoPage); err != nil {
			return da.EndWithError(err)
		}

		if err := getVideoPageMetadata(videoPage, videoId, rdx); err != nil {
			da.Error(err)
		}

		if err := downloadVideo(dolo.DefaultClient, videoId, videoPage, videoForce, videoSingleFormat); err != nil {
			da.Error(err)
		}

		if err := yeti.GetPosters(videoId, dolo.DefaultClient, force, yt_urls.AllThumbnailQualities()...); err != nil {
			da.Error(err)
		}

		if err := getVideoPageCaptions(videoPage, videoId, rdx, dolo.DefaultClient, force); err != nil {

			da.Error(err)
		}

		// set downloaded date
		if err := rdx.AddValues(data.VideoDownloadedDateProperty, videoId, time.Now().Format(time.RFC3339)); err != nil {
			return da.EndWithError(err)
		}

		// remove from the queue upon successful download
		if err := rdx.CutValues(data.VideosDownloadQueueProperty, videoId, data.TrueValue); err != nil {
			return da.EndWithError(err)
		}

		// add to watchlist upon successful download
		if err := rdx.AddValues(data.VideosWatchlistProperty, videoId, data.TrueValue); err != nil {
			return da.EndWithError(err)
		}

		da.Increment()
	}

	da.EndWithResult("done")

	return nil
}
