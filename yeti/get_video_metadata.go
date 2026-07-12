package yeti

import (
	"github.com/boggydigital/nod"
	"github.com/boggydigital/redux"
	"github.com/boggydigital/yet_urls/youtube_urls"
)

func GetVideoPageMetadata(videoPage *youtube_urls.InitialPlayerResponse, videoId string, rdx redux.Writeable) error {

	gvpma := nod.Begin(" metadata for %s", videoId)
	defer gvpma.Done()

	var err error
	if videoPage == nil {
		videoPage, err = GetVideoPage(videoId)
		if err != nil {
			return err
		}
	}

	for p, v := range ExtractMetadata(videoPage) {
		if err = rdx.AddValues(p, videoId, v...); err != nil {
			return err
		}
	}

	return nil
}
