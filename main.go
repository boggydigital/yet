package main

import (
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yt_urls"
	"log"
	"net/http"
	"os"
)

const (
	getVideosTopic = "getting video(s):"
)

var (
	streamingSources = []string{"best streaming quality", "good streaming quality", "available streaming quality"}
)

func main() {
	nod.EnableStdOut()

	if err := GetVideos(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func GetVideos(urlsOrVideoIds []string) error {
	nod.Start(getVideosTopic)

	dl := dolo.NewClient(http.DefaultClient, dolo.Defaults())

	for _, urlOrVideoId := range urlsOrVideoIds {

		videoId, err := yt_urls.VideoId(urlOrVideoId)
		if err != nil {
			return err
		}

		nod.Start(getVideosTopic, videoId)

		title, vidUrls, err := yt_urls.TitleStreamingUrls(videoId)
		if err != nil {
			return err
		}

		if len(vidUrls) == 0 {
			continue
		}

		attempt := 0
		for _, vidUrl := range vidUrls {

			topics := []string{getVideosTopic, videoId, streamingSources[attempt]}
			nod.Start(topics...)

			if vidUrl == nil || len(vidUrl.String()) == 0 {
				continue
			}

			tpw := nod.TotalProgress(topics...)

			_, err = dl.Download(vidUrl, "", saneFilename(title, videoId), tpw)

			if err != nil {
				attempt++
				if attempt > len(streamingSources)-1 {
					attempt = len(streamingSources) - 1
				}
				nod.Error(err, topics...)
				nod.End(topics...)
				continue
			}

			nod.End(getVideosTopic, videoId, streamingSources[attempt])
			nod.Success(true, topics...)

			//yt_urls.StreamingUrls returns bitrate sorted video urls,
			//so we can stop, if we've successfully got the best streaming quality
			break
		}

		nod.End(getVideosTopic, videoId)
	}

	nod.End(getVideosTopic)
	return nil
}
