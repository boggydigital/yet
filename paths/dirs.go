package paths

import (
	"github.com/boggydigital/pasu"
)

const DefaultYetRootDir = "/usr/share/yet"

const (
	Backups  pasu.AbsDir = "backups"
	Input    pasu.AbsDir = "input"
	Metadata pasu.AbsDir = "metadata"
	Videos   pasu.AbsDir = "videos"
	Posters  pasu.AbsDir = "posters"
	Captions pasu.AbsDir = "captions"
	Players  pasu.AbsDir = "players"
)

var AllAbsDirs = []pasu.AbsDir{
	Backups,
	Input,
	Metadata,
	Videos,
	Posters,
	Captions,
	Players,
}
