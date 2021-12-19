package main

import (
	"github.com/boggydigital/nod"
	"os"
)

func main() {
	nod.EnableStdOutPresenter()

	ya := nod.Begin("yetting requested videos/playlists")
	defer ya.End()

	//internally yet operates on video-ids, so the first step to process user input
	//is to expand all channel-ids into lists of video-ids and transparently return
	//any video-ids in the input stream
	videoIds, err := argsToVideoIds(os.Args[1:]...)
	if err != nil {
		_ = ya.EndWithError(err)
	}

	//having a list of video-ids, the only remaining thing is to download it one by one
	if err := DownloadVideos(videoIds...); err != nil {
		_ = ya.EndWithError(err)
	}
}
