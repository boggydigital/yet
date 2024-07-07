package yeti

import (
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yet_urls/youtube_urls"
	"net/url"
)

func GetCaptions(dl *dolo.Client, rdx kevlar.WriteableRedux, videoId string, captionTracks []youtube_urls.CaptionTrack, force bool) error {

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
			if err := dl.Download(u, force, nil, absFilename); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}
