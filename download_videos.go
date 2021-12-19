package main

import (
	"fmt"
	"github.com/boggydigital/cooja"
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yt_urls"
	"net/http"
	"net/url"
	"os"
)

func DownloadVideos(videoIds ...string) error {
	if len(videoIds) == 0 {
		return nil
	}

	dvtpw := nod.NewProgress(fmt.Sprintf("downloading %d video(s)", len(videoIds)))
	defer dvtpw.End()

	dvtpw.Total(uint64(len(videoIds)))

	jar, err := cooja.NewJar([]string{"youtube.com"}, "")
	if err != nil {
		return dvtpw.EndWithError(err)
	}

	dl := dolo.NewClient(&http.Client{Jar: jar}, dolo.Defaults())

	for _, videoId := range videoIds {

		gv := nod.Begin("video-id: " + videoId)

		vp, err := yt_urls.GetVideoPage(videoId)
		if err != nil {
			_ = gv.EndWithError(err)
			dvtpw.Increment()
			continue
		}

		title, vidUrls := vp.Title(), vp.StreamingFormats()

		if len(vidUrls) == 0 {
			_ = gv.EndWithError(err)
			dvtpw.Increment()
			continue
		}

		for _, vidUrl := range vidUrls {

			if vidUrl.Url == "" {
				continue
			}

			fn := saneFilename(title, videoId)
			tpw := nod.NewProgress("title: " + title)

			u, err := url.Parse(vidUrl.Url)
			if err != nil {
				_ = tpw.EndWithError(err)
				continue
			}

			if err := dl.Download(u, tpw, "", fn); err != nil {
				_ = tpw.EndWithError(err)
				continue
			}

			if _, err := os.Stat(fn); err == nil {
				if err := os.Chtimes(fn, vp.PublishDate(), vp.PublishDate()); err != nil {
					return tpw.EndWithError(err)
				}
			} else if os.IsNotExist(err) {
				//there was an error downloading this version, but we have a partial file
				//we can try resuming next application session
				break
			} else {
				return tpw.EndWithError(err)
			}

			tpw.EndWithResult("done")

			//yt_urls.StreamingUrls returns bitrate sorted video urls,
			//so we can stop, if we've successfully got the best streaming quality
			break
		}

		gv.End()
		dvtpw.Increment()
	}

	return nil
}
