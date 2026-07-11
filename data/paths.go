package data

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/boggydigital/camino"
	"github.com/boggydigital/yet_urls/youtube_urls"
)

const (
	defaultCaptionsExt = ".ytt"
	defaultScriptExt   = ".js"
)

// AbsPosterPath constructs poster path using poster directory,
// first and second letters of video-id, video-id itself
// and finally poster quality to get something like:
// /path/to/posters/v/i/videoId/quality.jpg
func AbsPosterPath(videoId string, quality youtube_urls.ThumbnailQuality) (string, error) {

	pdp := camino.GetAbs(Posters)

	spdp, err := mkdirAllVideoIdDirs(pdp, videoId)
	if err != nil {
		return "", err
	}

	return filepath.Join(spdp, quality.String()+youtube_urls.DefaultThumbnailExt), nil
}

// AbsCaptionsTrackPath constructs caption track path using captions directory,
// first and second letters of video-id to product something like
// /path/to/captions/f/s/fs_lang.jpg
func AbsCaptionsTrackPath(videoId, lang string) (string, error) {
	cdp := camino.GetAbs(Captions)

	scdp, err := mkdirAllVideoIdDirs(cdp, videoId)
	if err != nil {
		return "", err
	}

	return filepath.Join(scdp, fmt.Sprintf("%s_%s%s", videoId, lang, defaultCaptionsExt)), nil
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

	// finally, add videoId as the final path component
	subPath = filepath.Join(subPath, videoId)
	if _, err := os.Stat(subPath); os.IsNotExist(err) {
		if err := os.Mkdir(subPath, 777); err != nil {
			return "", err
		}
	}

	return subPath, nil
}
