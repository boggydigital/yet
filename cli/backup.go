package cli

import (
	"github.com/boggydigital/hogo"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pasu"
	"github.com/boggydigital/yet/paths"
	"net/url"
)

const daysToPreserveFiles = 30

func BackupHandler(_ *url.URL) error {
	return Backup()
}

func Backup() error {
	ea := nod.NewProgress("backing up metadata...")
	defer ea.End()

	amp, err := pasu.GetAbsDir(paths.Metadata)
	if err != nil {
		return ea.EndWithError(err)
	}

	abp, err := pasu.GetAbsDir(paths.Backups)
	if err != nil {
		return ea.EndWithError(err)
	}

	if err := hogo.Compress(amp, abp, ea); err != nil {
		return ea.EndWithError(err)
	}

	ea.EndWithResult("done")

	cba := nod.NewProgress("cleaning up old backups...")
	defer cba.End()

	if err := hogo.Cleanup(abp, true, cba); err != nil {
		return cba.EndWithError(err)
	}

	cba.EndWithResult("done")

	return nil
}
