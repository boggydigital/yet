package main

import (
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yt_urls"
	"os"
)

//GetPlaylistVideos returns all videoIds in a provided YouTube playlist,
//identified by a playlistId (the value of the "list" URL parameter).
//GetPlaylistVideos will skip videoId if there is an existing
//local file with a filename matching title, videoId.
//Note: GetPlaylistVideos can similarly enumerate videoIds for channels and users,
//given that (almost all) channel and user videos can be expressed as a
//playlist - e.g. "PLAY ALL" link for a channel/user "Videos" page is a playlist URL.
func GetPlaylistVideos(playlistId string) ([]string, error) {

	dp := nod.Start("downloading playlist: " + playlistId)

	playlist, err := yt_urls.GetPlaylistPage(playlistId)
	if err != nil {
		return nil, err
	}

	videoIds := make([]string, 0, len(playlist.Videos()))

	for playlist != nil &&
		len(playlist.Videos()) > 0 {
		for _, videoIdTitle := range playlist.Videos() {
			//before attempting to download - filter out the videos that are already present
			//locally, to download only updates to the playlist
			fn := saneFilename(videoIdTitle.Title, videoIdTitle.VideoId)
			if _, err := os.Stat(fn); err == nil {
				//file for the title, videoId combination has been downloaded already
				continue
			}

			videoIds = append(videoIds, videoIdTitle.VideoId)
		}

		if playlist.HasContinuation() {
			playlist, err = playlist.Continue()
			if err != nil {
				return videoIds, err
			}
		} else {
			playlist = nil
		}
	}

	dp.End()

	return videoIds, nil
}
