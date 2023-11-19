package yeti

import (
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yt_urls"
	"net/url"
)

func GetCaptions(dl *dolo.Client, rxa kvas.ReduxAssets, videoId string, captionTracks []yt_urls.CaptionTrack) error {

	if err := rxa.IsSupported(data.VideoCaptionsLanguages); err != nil {
		return err
	}

	captionLanguages := make([]string, 0, len(captionTracks))
	for _, ct := range captionTracks {
		captionLanguages = append(captionLanguages, ct.LanguageCode)
	}

	if err := rxa.ReplaceValues(data.VideoCaptionsLanguages, videoId, captionLanguages...); err != nil {
		return err
	}

	for _, ct := range captionTracks {

		u, err := url.Parse(ct.BaseUrl)
		if err != nil {
			return err
		}

		if absFilename, err := paths.AbsCaptionsTrackPath(videoId, ct.LanguageCode); err == nil {
			if err := dl.Download(u, nil, absFilename); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}
