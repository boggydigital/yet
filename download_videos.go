package main

import (
	"fmt"
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yt_urls"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func DownloadVideos(videoIds ...string) error {

	dv := nod.Start("downloading videos: " + strings.Join(videoIds, ", "))

	if len(videoIds) == 0 {
		return nod.Fatal(fmt.Errorf("you need to specify at least one video-id or URL"))
	}

	dl := dolo.NewClient(http.DefaultClient, dolo.Defaults())

	for _, videoId := range videoIds {

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

			topics := []string{videoIdTopic, titleTopic, "downloading..."}

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

			fn := saneFilename(title, videoId)

			_, err = dl.Download(u, "", fn, tpw)

			if err != nil {
				nod.Error(err, topics...)
				continue
			}

			if _, err := os.Stat(fn); err == nil {
				if err := os.Chtimes(fn, vp.PublishDate(), vp.PublishDate()); err != nil {
					return nod.Fatal(err, topics...)
				}
			} else if os.IsNotExist(err) {
				//there was an error downloading this version, but we have a partial file
				//we can try resuming next application session
				break
			} else {
				return nod.Fatal(err, topics...)
			}

			nod.EndResult("done", topics...)

			//yt_urls.StreamingUrls returns bitrate sorted video urls,
			//so we can stop, if we've successfully got the best streaming quality
			break
		}
	}

	dv.End()

	return nil
}
