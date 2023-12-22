package yeti

import (
	"github.com/boggydigital/coost"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yt_urls"
	"net/http"
	"strings"
)

func GetVideoPage(videoId string) (*yt_urls.InitialPlayerResponse, error) {

	// by default - use a default client that doesn't provide client cookies
	videoPage, err := yt_urls.GetVideoPage(http.DefaultClient, videoId)
	if err != nil {
		if strings.Contains(err.Error(), "Sign in to confirm your age") {
			// fallback to HTTP client with cookies
			absCookiePath, err := paths.AbsCookiesPath()
			if err != nil {
				return nil, err
			}
			if hc, err := coost.NewHttpClientFromFile(absCookiePath); err != nil {
				return nil, err
			} else {
				if videoPage, err = yt_urls.GetVideoPage(hc, videoId); err != nil {
					return nil, err
				}
			}
		} else {
			return nil, err
		}
	}
	return videoPage, nil
}
