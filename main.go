package main

import (
	"fmt"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yt_urls"
	"log"
	"os"
	"strings"
)

func main() {
	nod.EnableStdOut()

	videoIds, err := argsToVideoIds(os.Args[1:]...)
	if err != nil {
		log.Fatalln(err)
	}

	if err := DownloadVideos(videoIds...); err != nil {
		log.Fatalln(err)
	}
}

//argsToVideoIds converts list of videoIds, playlistIds, video URLs,
//playlist URLs (in any order and combination) to a list of videoIds.
//Inputs in unsupported format will result in error.
func argsToVideoIds(args ...string) ([]string, error) {
	videoIds := make([]string, 0)
	for _, urlOrId := range args {
		if len(urlOrId) < 12 {
			//currently, YouTube videoIds are exactly 11 characters,
			//meaning any URL containing videoId would be longer than 11 characters.
			videoIds = append(videoIds, urlOrId)
		} else if !strings.Contains(urlOrId, "?") {
			//currently, YouTube URLs would contain "?" query parameter separator,
			//meaning non-URL longer than 11 characters will be playlistId
			playlistVideoIds, err := GetPlaylistVideos(urlOrId)
			if err != nil {
				return videoIds, err
			}
			videoIds = append(videoIds, playlistVideoIds...)
		} else if !strings.Contains(urlOrId, "list=") {
			//currently, YouTube playlist URLs identify lists with a "list" parameter,
			//meaning that a URL (supported by yet) without it would be a video URL (/watch?v=videoId).
			videoId, err := yt_urls.VideoId(urlOrId)
			if err != nil {
				return videoIds, err
			}
			videoIds = append(videoIds, videoId)
		} else if strings.Contains(urlOrId, "v=") {
			//currently, YouTube video URLs identify videos with a "v" parameter.
			playlistId, err := yt_urls.PlaylistId(urlOrId)
			if err != nil {
				return videoIds, err
			}
			playlistVideoIds, err := GetPlaylistVideos(playlistId)
			if err != nil {
				return videoIds, err
			}
			videoIds = append(videoIds, playlistVideoIds...)
		} else {
			//provided input doesn't map to either:
			//-videoId: <12 characters long
			//-playlistId: >=12 characters long and doesn't contain query parameters separator
			//-playlist URL: URL containing query parameters separator and a "list" parameter
			//-video URL: URL containing query parameters separator and a "v" parameter (and doesn't contain "list" parameter)
			return videoIds, fmt.Errorf("unknown id or URL format: %s", urlOrId)
		}
	}
	return videoIds, nil
}
