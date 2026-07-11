package data

import (
	"os"

	"github.com/boggydigital/camino"
)

const (
	yetRootDir          = "/usr/share/yet"
	directoriesFilename = "directories.txt"
)

const (
	Backups camino.AbsDir = iota
	Input
	Metadata
	Videos
	Posters
	Captions
	YtDlp
)

var absDirNames = map[camino.AbsDir]string{
	Backups:  "backups",
	Input:    "input",
	Metadata: "metadata",
	Videos:   "videos",
	Posters:  "posters",
	Captions: "captions",
	YtDlp:    "yt-dlp",
}

const YtDlpPlugins camino.RelDir = iota

var relDirNames = map[camino.RelDir]string{
	YtDlpPlugins: "plugins",
}

var relAbsParents = map[camino.RelDir][]camino.AbsDir{
	YtDlpPlugins: {YtDlp},
}

func InitYetCamino() error {

	var overrides map[string]string

	if _, err := os.Stat(directoriesFilename); err == nil {
		if overrides, err = camino.ReadOverrides(directoriesFilename); err != nil {
			return err
		}
	}

	resolvedYetAbsPaths := camino.ResolveAbsPaths(yetRootDir, absDirNames, overrides)

	return camino.Register(resolvedYetAbsPaths, relDirNames, relAbsParents)
}
