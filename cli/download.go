package cli

import (
	"github.com/boggydigital/coost"
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
	queue := u.Query().Has("queue")
	force := u.Query().Has("force")
	return Download(ids, queue, force)
}

func Download(ids []string, queue, force bool) error {

	da := nod.NewProgress("downloading videos...")
	defer da.End()

	metadataDir, err := paths.GetAbsDir(paths.Metadata)
	if err != nil {
		return da.EndWithError(err)
	}

	absCookiePath, err := paths.AbsCookiesPath()
	if err != nil {
		return da.EndWithError(err)
	}

	rdx, err := kvas.ReduxWriter(metadataDir, data.AllProperties()...)
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

		videoPage, err := yt_urls.GetVideoPage(http.DefaultClient, videoId)
		if err != nil {
			if strings.Contains(err.Error(), "Sign in to confirm your age") {
				if hc, err := coost.NewHttpClientFromFile(absCookiePath); err != nil {
					return da.EndWithError(err)
				} else {
					if videoPage, err = yt_urls.GetVideoPage(hc, videoId); err != nil {
						return da.EndWithError(err)
					}
				}
			} else {
				return da.EndWithError(err)
			}
		}

		if err := getVideoPageMetadata(videoPage, videoId, rdx); err != nil {
			return da.EndWithError(err)
		}

		if err := downloadVideo(dolo.DefaultClient, videoId, videoPage); err != nil {
			return da.EndWithError(err)
		}

		if err := yeti.GetPosters(videoId, dolo.DefaultClient, yt_urls.AllThumbnailQualities()...); err != nil {
			return da.EndWithError(err)
		}

		if err := getVideoPageCaptions(videoPage, videoId, rdx, dolo.DefaultClient); err != nil {
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
