package cli

import (
	"github.com/boggydigital/coost"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yet/yeti"
	"net/url"
	"strings"
)

func DownloadHandler(u *url.URL) error {

	ids := strings.Split(u.Query().Get("id"), ",")
	later := u.Query().Has("later")
	return Download(ids, later)
}

func Download(ids []string, later bool) error {

	da := nod.Begin("downloading videos...")
	defer da.End()

	cookiesPath, err := paths.AbsCookiesPath()
	if err != nil {
		return da.EndWithError(err)
	}

	httpClient, err := coost.NewHttpClientFromFile(cookiesPath)
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

	if len(ids) > 0 {
		//internally yet operates on video-ids, so the first step to process user input
		//is to expand all channel-ids into lists of video-ids and transparently return
		//any video-ids in the input stream
		videoIds, err := yeti.ArgsToVideoIds(httpClient, false, ids...)
		if err != nil {
			return da.EndWithError(err)
		}

		// start by adding to download queue
		items, queueProperty := videoIds, data.VideosDownloadQueueProperty
		if len(videoIds) == 0 {
			items = ids
			queueProperty = data.UrlsDownloadQueueProperty
		}
		for _, item := range items {
			if err := rxa.AddValues(queueProperty, item, data.TrueValue); err != nil {
				return da.EndWithError(err)
			}
		}

		if later {
			return nil
		}

		if len(videoIds) > 0 {
			//having a list of video-ids, the only remaining thing is to download it one by one
			if err := yeti.DownloadVideos(httpClient, rxa, videoIds...); err != nil {
				return da.EndWithError(err)
			}
		} else {
			//argument has not been determined to be a video-id, attempt direct URL download
			if err := yeti.DownloadUrls(httpClient, rxa, ids...); err != nil {
				return da.EndWithError(err)
			}
		}

	} else {
		da.EndWithResult("expected one or more video-id, playlist-id")
	}
	return nil
}
