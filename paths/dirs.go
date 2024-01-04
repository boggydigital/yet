package paths

import (
	"github.com/boggydigital/pathology"
)

const DefaultYetRootDir = "/usr/share/yet"

const (
	Backups  pathology.AbsDir = "backups"
	Input    pathology.AbsDir = "input"
	Metadata pathology.AbsDir = "metadata"
	Videos   pathology.AbsDir = "videos"
	Posters  pathology.AbsDir = "posters"
	Captions pathology.AbsDir = "captions"
)

var AllAbsDirs = []pathology.AbsDir{
	Backups,
	Input,
	Metadata,
	Videos,
	Posters,
	Captions,
}
