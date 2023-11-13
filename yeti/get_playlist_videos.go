package yeti

import (
	"fmt"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yt_urls"
	"net/http"
	"os"
)

// GetPlaylistVideos returns all videoIds in a provided YouTube playlist,
// identified by a playlistId (the value of the "list" URL parameter).
// GetPlaylistVideos will skip videoId if there is an existing
// local file with a filename matching title, videoId.
// Note: GetPlaylistVideos can similarly enumerate videoIds for channels and users,
// given that (almost all) channel and user videos can be expressed as a
// playlist - e.g. "PLAY ALL" link for a channel/user "Videos" page is a playlist URL.
func GetPlaylistVideos(httpClient *http.Client, playlistId string, newVideos bool) ([]string, error) {

	dp := nod.Begin(fmt.Sprintf("itemizing playlist %s:", playlistId))
	defer dp.End()

	playlistHasVideos := false

	playlist, err := yt_urls.GetPlaylistPage(httpClient, playlistId)
	if err != nil {
		return nil, dp.EndWithError(err)
	}

	videoIds := make([]string, 0, len(playlist.Videos()))

	for playlist != nil &&
		len(playlist.Videos()) > 0 {
		playlistHasVideos = true
		for _, vtc := range playlist.Videos() {
			//before attempting to download - filter out the videos that are already present
			//locally, to download only updates to the playlist
			fn := ChannelTitleVideoIdFilename(vtc.Channel, vtc.Title, vtc.VideoId)
			if _, err := os.Stat(fn); err == nil {
				//file for the title, videoId combination has been downloaded already
				//TODO: change from returning on first existing video to return on no new videos in playlist page
				if newVideos {
					if len(videoIds) == 0 {
						dp.EndWithResult("no new videos")
					} else {
						dp.EndWithResult("found %d new video(s)", len(videoIds))
					}
					return videoIds, nil
				} else {
					continue
				}
			}

			videoIds = append(videoIds, vtc.VideoId)
		}

		if playlist.HasContinuation() {
			playlist, err = playlist.Continue(httpClient)
			if err != nil {
				return videoIds, dp.EndWithError(err)
			}
		} else {
			playlist = nil
		}
	}

	if len(videoIds) > 0 {
		dp.EndWithResult("got %d video(s)", len(videoIds))
	} else if playlistHasVideos {
		dp.EndWithResult("no new videos to download")
	}

	return videoIds, nil
}
