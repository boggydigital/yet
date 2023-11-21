package yeti

import (
	"fmt"
	"github.com/boggydigital/yt_urls"
	"strings"
)

const (
	youtuBeHost = "youtu.be/"
)

// ParseVideoIds converts list of videoIds in any form - as video-ids,
// YouTube /watch, youtu.be/ URLs (in any order and combination) to a list of videoIds.
// Inputs in unsupported format will produce an error.
func ParseVideoIds(args ...string) ([]string, error) {
	videoIds := make([]string, 0)
	for _, urlOrId := range args {
		if len(urlOrId) < 12 {
			//currently, YouTube videoIds are exactly 11 characters,
			//meaning any URL containing videoId would be longer than 11 characters.
			videoIds = append(videoIds, urlOrId)
		} else if strings.Contains(urlOrId, youtuBeHost) {
			//currently, YouTube own short URLs are formatted as
			//youtu.be/videoId
			if _, videoId, ok := strings.Cut(urlOrId, youtuBeHost); ok {
				videoIds = append(videoIds, videoId)
			}
		} else if strings.Contains(urlOrId, "v=") {
			//currently, YouTube video URLs contain v=video-id parameter
			videoId, err := yt_urls.VideoId(urlOrId)
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