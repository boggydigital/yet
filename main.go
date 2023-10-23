package main

import (
	"github.com/boggydigital/coost"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/yeti"
	"os"
)

func main() {
	nod.EnableStdOutPresenter()

	ya := nod.Begin("yet is getting requested videos/playlists")
	defer ya.End()

	bins := yeti.NewBinaries()

	httpClient, err := coost.NewHttpClientFromFile("cookies.txt")
	if err != nil {
		_ = ya.EndWithError(err)
		os.Exit(1)
	}

	if len(os.Args) > 1 {
		//internally yet operates on video-ids, so the first step to process user input
		//is to expand all channel-ids into lists of video-ids and transparently return
		//any video-ids in the input stream
		videoIds, err := yeti.ArgsToVideoIds(httpClient, false, os.Args[1:]...)
		if err != nil {
			_ = ya.EndWithError(err)
		}

		if len(videoIds) > 0 {
			//having a list of video-ids, the only remaining thing is to download it one by one
			if err := yeti.DownloadVideos(httpClient, yeti.DefaultFilenameDelegate, bins, videoIds...); err != nil {
				_ = ya.EndWithError(err)
			}
		} else {
			//argument has not been determined to be a video-id, attempt direct URL download
			if err := yeti.DownloadUrls(httpClient, os.Args[1:]...); err != nil {
				_ = ya.EndWithError(err)
			}
		}

		return
	}

	ya.EndWithResult("...or not:")
	nod.ErrorStr("No arguments specified, expected: yet <video-id>[, ...] <channel-id>[, ...] ")
}
