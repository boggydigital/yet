package main

import (
	"fmt"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yt_urls"
)

func GetPlaylistVideos(playlistId string) error {

	dp := nod.Start("downloading playlist: " + playlistId)

	playlist, err := yt_urls.GetPlaylistPage(playlistId)
	if err != nil {
		return err
	}

	for playlist != nil {
		for _, videoIdTitle := range playlist.Videos() {
			if err := GetVideos(videoIdTitle.VideoId); err != nil {
				fmt.Println(err)
			}
		}
		if playlist.HasContinuation() {
			playlist, err = playlist.Continue()
			if err != nil {
				return err
			}
		} else {
			playlist = nil
		}
	}

	dp.End()

	return nil
}
