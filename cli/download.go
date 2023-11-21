package cli

import (
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/yeti"
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

	if len(ids) > 0 {
		//internally yet operates on video-ids, so the first step to process user input
		//is to expand all channel-ids into lists of video-ids and transparently return
		//any video-ids in the input stream
		videoIds, err := yeti.ParseVideoIds(ids...)
		if err != nil {
			return da.EndWithError(err)
		}

		if len(videoIds) > 0 {
			//having a list of video-ids, the only remaining thing is to download it one by one
			if err := GetVideo(force, videoIds...); err != nil {
				return da.EndWithError(err)
			}
		} else {
			//argument has not been determined to be a video-id, attempt direct URL download
			if err := GetFile(ids...); err != nil {
				return da.EndWithError(err)
			}
		}

	} else {
		da.EndWithResult("expected one or more video-id, playlist-id")
	}
	return nil
}
