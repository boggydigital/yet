package main

import (
	"fmt"
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yt_urls"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	nod.EnableStdOut()

	if err := GetVideos(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func GetVideos(urlsOrVideoIds []string) error {

	if len(urlsOrVideoIds) == 0 {
		return fmt.Errorf("you need to specify at least one video-id or URL")
	}

	dl := dolo.NewClient(http.DefaultClient, dolo.Defaults())

	for _, urlOrVideoId := range urlsOrVideoIds {

		videoId, err := yt_urls.VideoId(urlOrVideoId)
		if err != nil {
			return err
		}

		videoIdTopic := "getting video-id: " + videoId
		nod.Start(videoIdTopic)

		vp, err := yt_urls.GetVideoPage(videoId)
		if err != nil {
			return err
		}

		title, vidUrls := vp.Title(), vp.StreamingFormats()

		titleTopic := "title: " + title
		if title != "" {
			nod.Start(videoIdTopic, titleTopic)
		}

		if len(vidUrls) == 0 {
			continue
		}

		for _, vidUrl := range vidUrls {

			qualityTopic := fmt.Sprintf("downloading...")
			topics := []string{videoIdTopic, titleTopic, qualityTopic}

			nod.Start(topics...)

			if vidUrl.Url == "" {
				continue
			}

			tpw := nod.TotalProgress(topics...)

			u, err := url.Parse(vidUrl.Url)
			if err != nil {
				nod.Error(err, topics...)
				continue
			}

			_, err = dl.Download(u, "", saneFilename(title, videoId), tpw)

			if err != nil {
				nod.Error(err, topics...)
				continue
			}

			nod.Result("done", topics...)

			//yt_urls.StreamingUrls returns bitrate sorted video urls,
			//so we can stop, if we've successfully got the best streaming quality
			break
		}
	}

	return nil
}

func saneFilename(title, videoId string) string {
	unsafeChars := []string{"/", ":", "?", "*"}

	fn := fmt.Sprintf("%s-%s", title, videoId)
	if title == "" {
		fn = fmt.Sprintf("%s", videoId)
	}

	for _, ch := range unsafeChars {
		fn = strings.ReplaceAll(fn, ch, "")
	}

	fn += ".mp4"

	return fn
}
