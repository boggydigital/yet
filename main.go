package main

import (
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pesco"
	"github.com/boggydigital/yt_urls"
	"os"
)

func main() {
	nod.EnableStdOutPresenter()

	ya := nod.Begin("yetting requested videos/playlists")
	defer ya.End()

	jar, err := pesco.NewJar([]string{yt_urls.YoutubeHost}, "")
	if err != nil {
		_ = ya.EndWithError(err)
	}

	defer func(jar pesco.PersistentCookieJar) {
		if err := jar.Store(); err != nil {
			_ = ya.EndWithError(err)
		}
	}(jar)

	httpClient := jar.NewClient()

	//internally yet operates on video-ids, so the first step to process user input
	//is to expand all channel-ids into lists of video-ids and transparently return
	//any video-ids in the input stream
	videoIds, err := argsToVideoIds(httpClient, os.Args[1:]...)
	if err != nil {
		_ = ya.EndWithError(err)
	}

	//having a list of video-ids, the only remaining thing is to download it one by one
	if err := DownloadVideos(httpClient, videoIds...); err != nil {
		_ = ya.EndWithError(err)
	}
}
