package yeti

import (
	"net/http"
	"net/url"
	"strings"
)

func PlaylistVideoIds(httpClient *http.Client, newPlaylistVideos bool, args ...string) ([]string, error) {
	videoIds := make([]string, 0)
	for _, urlOrId := range args {
		if strings.Contains(urlOrId, "list=") {
			//currently, YouTube playlist URLs would contain "/playlist" endpoint
			u, err := url.Parse(urlOrId)
			if err != nil {
				return nil, err
			}
			listId := u.Query().Get("list")
			playlistVideoIds, err := GetPlaylistVideos(httpClient, listId, newPlaylistVideos)
			if err != nil {
				return videoIds, err
			}
			videoIds = append(videoIds, playlistVideoIds...)
		}
	}

	return videoIds, nil
}
