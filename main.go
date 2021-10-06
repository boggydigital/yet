package main

import (
	"fmt"
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yt_urls"
	"net/http"
	"os"
	"strings"
)

func main() {
	nod.EnableStdOut()

	if err := GetVideos(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func GetVideos(urlsOrVideoIds []string) error {
	//nod.Begin(getVideosTopic)

	dl := dolo.NewClient(http.DefaultClient, nil, dolo.Defaults())

	for _, urlOrVideoId := range urlsOrVideoIds {

		videoId, err := yt_urls.VideoId(urlOrVideoId)
		if err != nil {
			return err
		}

		//nod.Begin(getVideosTopic, videoId)

		title, vidUrls, err := yt_urls.TitleStreamingUrls(videoId)
		if err != nil {
			return err
		}

		if len(vidUrls) == 0 {
			continue
		}

		for _, vidUrl := range vidUrls {

			if vidUrl == nil || len(vidUrl.String()) == 0 {
				continue
			}

			filename := fmt.Sprintf("%s-%s", title, videoId)
			if title == "" {
				filename = fmt.Sprintf("%s", videoId)
			}

			filename = strings.ReplaceAll(filename, "/", "")
			filename = strings.ReplaceAll(filename, ":", "")
			filename = strings.ReplaceAll(filename, ".", "")

			filename += ".mp4"

			_, err = dl.Download(vidUrl, "./", filename)
			if err != nil {
				fmt.Println(err)
				continue
			}

			//yt_urls.StreamingUrls returns bitrate sorted video urls,
			//so we can stop, if we've successfully got the best available one
			break
		}

		//nod.End(getVideosTopic, videoId)
	}

	//nod.End(getVideosTopic)
	return nil
}
