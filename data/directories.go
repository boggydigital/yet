package data

import (
	"github.com/boggydigital/pathways"
)

const DefaultRootDir = "/usr/share/yet"

const (
	Backups  pathways.AbsDir = "backups"
	Input    pathways.AbsDir = "input"
	Metadata pathways.AbsDir = "metadata"
	Videos   pathways.AbsDir = "videos"
	Posters  pathways.AbsDir = "posters"
	Captions pathways.AbsDir = "captions"
	YtDlp    pathways.AbsDir = "yt-dlp"
)

const (
	YtDlpPlugins pathways.RelDir = "plugins"
)

var RelToAbsDirs = map[pathways.RelDir]pathways.AbsDir{
	YtDlpPlugins: YtDlp,
}

var AllAbsDirs = []pathways.AbsDir{
	Backups,
	Input,
	Metadata,
	Videos,
	Posters,
	Captions,
	YtDlp,
}
