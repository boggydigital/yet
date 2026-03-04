package yeti

import (
	"net/http"

	"github.com/boggydigital/yet_urls/youtube_urls"
)

var errorsSolvedWithCookies = []string{
	"Sign in to confirm your age",
	"Join this channel to get access to members-only content",
}

func GetVideoPage(videoId string) (*youtube_urls.InitialPlayerResponse, error) {

	// by default - use a default client that doesn't supply client cookies
	videoPage, err := youtube_urls.GetVideoPage(http.DefaultClient, videoId)
	if err != nil {
		// TODO: rewrite this to handle more gracefully
		//errSolvedWithCookies := false
		//for _, esc := range errorsSolvedWithCookies {
		//	if strings.Contains(err.Error(), esc) {
		//		errSolvedWithCookies = true
		//		// fallback to HTTP client with cookies
		//		absCookiePath := data.AbsCookiesPath()
		//		var jar http.CookieJar
		//		jar, err = coost.Read(youtube_urls.HostUrl(), absCookiePath)
		//		if err != nil {
		//			return nil, err
		//		}
		//		client := http.DefaultClient
		//		client.Jar = jar
		//
		//		return youtube_urls.GetVideoPage(client, videoId)
		//	} else {
		//		return nil, err
		//	}
		//}
		//if !errSolvedWithCookies {
		//	return nil, err
		//}
		return nil, err
	}
	return videoPage, nil
}
