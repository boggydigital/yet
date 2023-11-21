package paths

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	cookiesFilename      = "cookies.txt"
	defaultPosterExt     = ".jpg"
	defaultYTCaptionsExt = ".ytt"
	PosterQualityMax     = "maxresdefault"
	PosterQualityHigh    = "hqdefault"
)

func AbsCookiesPath() (string, error) {
	idp, err := GetAbsDir(Input)
	return filepath.Join(idp, cookiesFilename), err
}

// AbsPosterPath constructs poster path using poster directory,
// first and second letters of video-id to product something like
// /path/to/posters/f/s/fs_quality.jpg
func AbsPosterPath(videoId, quality string) (string, error) {

	pdp, err := GetAbsDir(Posters)
	if err != nil {
		return "", err
	}

	spdp, err := mkdirAllVideoIdDirs(pdp, videoId)
	if err != nil {
		return "", err
	}

	return filepath.Join(spdp, fmt.Sprintf("%s_%s%s", videoId, quality, defaultPosterExt)), nil
}

// AbsCaptionsTrackPath constructs caption track path using captions directory,
// first and second letters of video-id to product something like
// /path/to/captions/f/s/fs_lang.jpg
func AbsCaptionsTrackPath(videoId, lang string) (string, error) {
	cdp, err := GetAbsDir(Captions)
	if err != nil {
		return "", err
	}

	scdp, err := mkdirAllVideoIdDirs(cdp, videoId)
	if err != nil {
		return "", err
	}

	return filepath.Join(scdp, fmt.Sprintf("%s_%s%s", videoId, lang, defaultYTCaptionsExt)), nil
}

func mkdirAllVideoIdDirs(path, videoId string) (string, error) {

	if len(videoId) < 2 {
		return "", errors.New("video-id is too short to construct sub-path")
	}

	// add the first video-id letter to the sub-path
	subPath := filepath.Join(path, videoId[:1])
	if _, err := os.Stat(subPath); os.IsNotExist(err) {
		if err := os.Mkdir(subPath, 777); err != nil {
			return "", err
		}
	}

	// add the second video-id letter to the sub-path
	subPath = filepath.Join(subPath, videoId[1:2])
	if _, err := os.Stat(subPath); os.IsNotExist(err) {
		if err := os.Mkdir(subPath, 777); err != nil {
			return "", err
		}
	}

	return subPath, nil
}
