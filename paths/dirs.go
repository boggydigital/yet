package paths

import (
	"fmt"
	"os"
)

type AbsDir int

const (
	Backups AbsDir = iota
	Input
	Metadata
	Videos
	Posters
	Captions
)

var absDirsStrings = map[AbsDir]string{
	Backups:  "backups",
	Input:    "input",
	Metadata: "metadata",
	Videos:   "videos",
	Posters:  "posters",
	Captions: "captions",
}

var absDirsPaths = map[AbsDir]string{}

func SetAbsDirs(kv map[string]string) error {
	for adk, ads := range absDirsStrings {
		if d, ok := kv[ads]; ok && d != "" {
			if d != DefaultDirs[ads] {
				// make sure directory exists
				if _, err := os.Stat(d); err != nil {
					return err
				}
			}
			absDirsPaths[adk] = d
		} else {
			return fmt.Errorf("missing required abs dir %s", ads)
		}
	}
	return nil
}

func GetAbsDir(ad AbsDir) (string, error) {
	if _, ok := absDirsStrings[ad]; !ok {
		return "", fmt.Errorf("unknown abs dir")
	}

	if adp, ok := absDirsPaths[ad]; ok && adp != "" {
		return adp, nil
	}
	return "", fmt.Errorf("abs dir %s not set", absDirsStrings[ad])
}
