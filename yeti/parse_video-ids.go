package yeti

import (
	"fmt"
	"github.com/boggydigital/yet_urls/youtube_urls"
	"net/url"
	"path"
	"strings"
)

const (
	youtuBeHost     = "youtu.be/"
	youtubeEmbedUrl = "youtube.com/embed"
)

// ParseVideoIds converts list of videoIds in any form - as video-ids,
// YouTube /watch, youtu.be/ URLs (in any order and combination) to a list of videoIds.
// Inputs in unsupported format will produce an error.
func ParseVideoIds(args ...string) ([]string, error) {
	videoIds := make([]string, 0)
	for _, urlOrId := range args {
		if urlOrId == "" {
			continue
		} else if len(urlOrId) < 12 {
			//currently, YouTube videoIds are exactly 11 characters,
			//meaning any URL containing videoId would be longer than 11 characters.
			videoIds = append(videoIds, urlOrId)
		} else if strings.Contains(urlOrId, youtuBeHost) {
			//currently, YouTube own short URLs are formatted as
			//youtu.be/videoId
			if ybeu, err := url.Parse(urlOrId); err == nil {
				videoIds = append(videoIds, path.Base(ybeu.Path))
			}
		} else if strings.Contains(urlOrId, youtubeEmbedUrl) {
			//currently YouTube embed links are formatted as
			//youtube.com/embed/videoId
			if u, err := url.Parse(urlOrId); err == nil {
				videoIds = append(videoIds, path.Base(u.Path))
			} else {
				return nil, err
			}
		} else if strings.Contains(urlOrId, "v=") {
			//currently, YouTube video URLs contain v=video-id parameter
			videoId, err := youtube_urls.VideoId(urlOrId)
			if err != nil {
				return videoIds, err
			}
			videoIds = append(videoIds, videoId)
		} else {
			//provided input doesn't map to either:
			//-videoId: <12 characters long
			//-youtu.be/videoId
			//-video URL: URL containing a "v=videoId" parameter
			//
			//that's currently not supported as a videoId input
			return nil, fmt.Errorf("%s is not a valid video-id input", urlOrId)
		}
	}
	return videoIds, nil
}

func ParseVideoId(videoId string) (string, error) {
	parsedVideoIds, err := ParseVideoIds(videoId)
	if err != nil {
		return "", err
	}
	if len(parsedVideoIds) > 0 {
		return parsedVideoIds[0], nil
	} else {
		return "", fmt.Errorf("invalid video id: %s", videoId)
	}
}
