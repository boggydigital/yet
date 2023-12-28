package cli

import (
	"net/url"
)

func BackupHandler(u *url.URL) error {
	return Backup()
}

func Backup() error {
	return nil
}
