package main

import (
	"github.com/boggydigital/coost"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yt_urls"
	"os"
)

func main() {
	nod.EnableStdOutPresenter()

	ya := nod.Begin("yet is getting requested videos/playlists")
	defer ya.End()

	jar, err := coost.NewJar([]string{yt_urls.YoutubeHost}, "")
	if err != nil {
		_ = ya.EndWithError(err)
	}

	httpClient := jar.NewHttpClient()

	if len(os.Args) > 1 {
		//internally yet operates on video-ids, so the first step to process user input
		//is to expand all channel-ids into lists of video-ids and transparently return
		//any video-ids in the input stream
		videoIds, err := argsToVideoIds(httpClient, false, os.Args[1:]...)
		if err != nil {
			_ = ya.EndWithError(err)
		}

		//having a list of video-ids, the only remaining thing is to download it one by one
		if err := DownloadVideos(httpClient, videoIds...); err != nil {
			_ = ya.EndWithError(err)
		}
	} else {
		ya.EndWithResult("...or not:")
		ha := nod.Begin("No arguments specified, expected: yet <video-id>[, ...] <channel-id>[, ...] ")
		ha.End()
	}
}
