package paths

import "path/filepath"

const (
	cookiesFilename = "cookies.txt"
)

func AbsCookiesPath() (string, error) {
	idp, err := GetAbsDir(Input)
	return filepath.Join(idp, cookiesFilename), err
}
