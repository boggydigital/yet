package yeti

import (
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yt_urls"
	"net/url"
)

func GetCaptions(dl *dolo.Client, rdx kvas.WriteableRedux, videoId string, captionTracks []yt_urls.CaptionTrack) error {

	properties := []string{
		data.VideoCaptionsNamesProperty,
		data.VideoCaptionsKindsProperty,
		data.VideoCaptionsLanguagesProperty}

	if err := rdx.MustHave(properties...); err != nil {
		return err
	}

	captionsData := make(map[string][]string)
	for _, p := range properties {
		captionsData[p] = make([]string, 0, len(captionTracks))
		for _, ct := range captionTracks {
			value := ""
			switch p {
			case data.VideoCaptionsNamesProperty:
				value = ct.TrackName
			case data.VideoCaptionsKindsProperty:
				value = ct.Kind
			case data.VideoCaptionsLanguagesProperty:
				value = ct.LanguageCode
			}
			captionsData[p] = append(captionsData[p], value)
		}
		if err := rdx.ReplaceValues(p, videoId, captionsData[p]...); err != nil {
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
