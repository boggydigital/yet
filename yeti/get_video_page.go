package yeti

import (
	"github.com/boggydigital/coost"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yet_urls/youtube_urls"
	"net/http"
	"strings"
)

var errorsSolvedWithCookies = []string{
	"Sign in to confirm your age",
	"Join this channel to get access to members-only content",
}

func GetVideoPage(videoId string) (*youtube_urls.InitialPlayerResponse, error) {

	// by default - use a default client that doesn't supply client cookies
	videoPage, err := youtube_urls.GetVideoPage(http.DefaultClient, videoId)
	if err != nil {
		errSolvedWithCookies := false
		for _, esc := range errorsSolvedWithCookies {
			if strings.Contains(err.Error(), esc) {
				errSolvedWithCookies = true
				// fallback to HTTP client with cookies
				if absCookiePath, err := paths.AbsCookiesPath(); err == nil {
					if hc, err := coost.NewHttpClientFromFile(absCookiePath); err == nil {
						return youtube_urls.GetVideoPage(hc, videoId)
					} else {
						return nil, err
					}
				} else {
					return nil, err
				}
			}
		}
		if !errSolvedWithCookies {
			return nil, err
		}
	}
	return videoPage, nil
}
