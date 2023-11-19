package paths

import (
	"fmt"
	"path/filepath"
)

const (
	cookiesFilename    = "cookies.txt"
	defaultPosterExt   = ".jpg"
	defaultCaptionsExt = ".vtt"
	PosterQualityMax   = "maxresdefault"
	PosterQualityHigh  = "hqdefault"
)

func AbsCookiesPath() (string, error) {
	idp, err := GetAbsDir(Input)
	return filepath.Join(idp, cookiesFilename), err
}

func AbsPosterPath(videoId, quality string) (string, error) {
	pdp, err := GetAbsDir(Posters)
	return filepath.Join(pdp, fmt.Sprintf("%s_%s%s", videoId, quality, defaultPosterExt)), err
}

func AbsCaptionsTrackPath(videoId, lang string) (string, error) {
	cdp, err := GetAbsDir(Captions)
	return filepath.Join(cdp, fmt.Sprintf("%s_%s%s", videoId, lang, defaultCaptionsExt)), err
}
