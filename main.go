package main

import (
	"github.com/boggydigital/coost"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yt_urls"
	"os"
	"os/exec"
)

const ffmpegCmdEnv = "YET_FFMPEG_CMD"

func main() {
	nod.EnableStdOutPresenter()

	ya := nod.Begin("yet is getting requested videos/playlists")
	defer ya.End()

	//get ffmpeg binary location from user specified env or elsewhere on the system
	ffmpegCmd := os.Getenv(ffmpegCmdEnv)
	if ffmpegCmd == "" {
		if path, err := exec.LookPath("ffmpeg"); err == nil {
			ffmpegCmd = path
		}
	}

	httpClient, err := coost.NewHttpClientFromFile("cookies.txt", yt_urls.YoutubeHost)
	if err != nil {
		_ = ya.EndWithError(err)
		os.Exit(1)
	}

	if len(os.Args) > 1 {
		//internally yet operates on video-ids, so the first step to process user input
		//is to expand all channel-ids into lists of video-ids and transparently return
		//any video-ids in the input stream
		videoIds, err := argsToVideoIds(httpClient, false, os.Args[1:]...)
		if err != nil {
			_ = ya.EndWithError(err)
		}

		//having a list of video-ids, the only remaining thing is to download it one by one
		if err := DownloadVideos(httpClient, localVideoFilename, ffmpegCmd, videoIds...); err != nil {
			_ = ya.EndWithError(err)
		}

		return
	}

	ya.EndWithResult("...or not:")
	nod.ErrorStr("No arguments specified, expected: yet <video-id>[, ...] <channel-id>[, ...] ")
}
