package yeti

import (
	"fmt"
	"github.com/boggydigital/yt_urls"
	"strings"
)

// ParsePlaylistIds converts list of playlistIds in any form - as playlist-ids,
// YouTube /watch, youtu.be/ URLs (in any order and combination) to a list of videoIds.
// Inputs in unsupported format will produce an error.
func ParsePlaylistIds(args ...string) ([]string, error) {
	playlistIds := make([]string, 0)
	for _, urlOrId := range args {
		if urlOrId == "" {
			continue
		} else if strings.HasPrefix(urlOrId, "PL") {
			//currently, playlists commonly have "PL...." prefix
			playlistIds = append(playlistIds, urlOrId)
		} else if strings.HasPrefix(urlOrId, "UU") {
			//currently, autogenerated playlists commonly have "UU..." prefix
			playlistIds = append(playlistIds, urlOrId)
		} else if strings.Contains(urlOrId, "list=") {
			//currently, YouTube playlist URLs contain list=playlist-id parameter
			playlistId, err := yt_urls.PlaylistId(urlOrId)
			if err != nil {
				return playlistIds, err
			}
			playlistIds = append(playlistIds, playlistId)
		} else {
			//provided input doesn't map to either:
			//-playlistId starting with "PL" or "UU"
			//-playlist URL: URL containing a "list=playlistId" parameter
			//
			//that's currently not supported as a playlistId input
			return nil, fmt.Errorf("%s is not a valid playlist-id input", urlOrId)
		}
	}
	return playlistIds, nil
}
