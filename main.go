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

	//defer func(jar coost.PersistentCookieJar) {
	//	if err := jar.Store(); err != nil {
	//		_ = ya.EndWithError(err)
	//	}
	//}(jar)

	httpClient := jar.NewHttpClient()

	args := os.Args[1:]

	if len(args) > 0 {
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
		//check if yet-list.txt is present and download specified videos
		//or print help if yet-list.txt is not present or is empty

		dirIds, err := readList()
		if err != nil {
			_ = ya.EndWithError(err)
		}

		if len(dirIds) == 0 {
			printHelp()
			return
		}

		if err := processList(dirIds); err != nil {
			_ = ya.EndWithError(err)
		}
	}
}
