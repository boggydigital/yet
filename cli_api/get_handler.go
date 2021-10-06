package cli_api

import (
	"fmt"
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/yt_urls"
	"net/http"
	"net/url"
	"strings"
)

const (
	getVideosTopic = "getting videos:"
)

func GetHandler(u *url.URL) error {
	q := u.Query()
	urlsOrVideoIdsStr := q.Get("video-id")
	urlsOrVideoIds := strings.Split(urlsOrVideoIdsStr, ",")
	return GetVideos(urlsOrVideoIds)
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
