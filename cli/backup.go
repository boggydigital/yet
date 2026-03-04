package cli

import (
	"net/url"

	"github.com/boggydigital/backups"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
)

func BackupHandler(_ *url.URL) error {
	return Backup()
}

func Backup() error {
	ea := nod.NewProgress("backing up metadata...")
	defer ea.Done()

	amp := data.Pwd.AbsDirPath(data.Metadata)
	abp := data.Pwd.AbsDirPath(data.Backups)

	if err := backups.Compress(amp, abp); err != nil {
		return err
	}

	cba := nod.NewProgress("cleaning up old backups...")
	defer cba.Done()

	if err := backups.Cleanup(abp, true, cba); err != nil {
		return err
	}

	return nil
}
