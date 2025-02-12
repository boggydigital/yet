package cli

import (
	"github.com/boggydigital/backups"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pathways"
	"github.com/boggydigital/yet/data"
	"net/url"
)

func BackupHandler(_ *url.URL) error {
	return Backup()
}

func Backup() error {
	ea := nod.NewProgress("backing up metadata...")
	defer ea.Done()

	amp, err := pathways.GetAbsDir(data.Metadata)
	if err != nil {
		return err
	}

	abp, err := pathways.GetAbsDir(data.Backups)
	if err != nil {
		return err
	}

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
