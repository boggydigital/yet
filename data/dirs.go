package data

import (
	"github.com/boggydigital/pathways"
)

const DefaultYetRootDir = "/usr/share/yet"

const (
	Backups  pathways.AbsDir = "backups"
	Input    pathways.AbsDir = "input"
	Metadata pathways.AbsDir = "metadata"
	Videos   pathways.AbsDir = "videos"
	Posters  pathways.AbsDir = "posters"
	Captions pathways.AbsDir = "captions"
	Players  pathways.AbsDir = "players"
	YtDlp    pathways.AbsDir = "yt-dlp"
)

var AllAbsDirs = []pathways.AbsDir{
	Backups,
	Input,
	Metadata,
	Videos,
	Posters,
	Captions,
	Players,
	YtDlp,
}
