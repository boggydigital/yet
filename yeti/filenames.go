package yeti

import (
	"path/filepath"
	"strings"
)

const (
	videoExt = ".video"
	audioExt = ".audio"
)

func videoAudioFilenames(relFilename string) (string, string) {
	ext := filepath.Ext(relFilename)
	fse := strings.TrimSuffix(relFilename, ext)

	return fse + videoExt, fse + audioExt
}
