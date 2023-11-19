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

	properties := []string{
		data.VideoCaptionsNames,
		data.VideoCaptionsKinds,
		data.VideoCaptionsLanguages}

	if err := rxa.IsSupported(properties...); err != nil {
		return err
	}

	captionsData := make(map[string][]string)
	for _, p := range properties {
		captionsData[p] = make([]string, 0, len(captionTracks))
		for _, ct := range captionTracks {
			value := ""
			switch p {
			case data.VideoCaptionsNames:
				value = ct.TrackName
			case data.VideoCaptionsKinds:
				value = ct.Kind
			case data.VideoCaptionsLanguages:
				value = ct.LanguageCode
			}
			captionsData[p] = append(captionsData[p], value)
		}
		if err := rxa.ReplaceValues(p, videoId, captionsData[p]...); err != nil {
			return err
		}
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