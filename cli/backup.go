package cli

import (
	"github.com/boggydigital/nod"
	"github.com/boggydigital/packer"
	"github.com/boggydigital/pathology"
	"github.com/boggydigital/yet/paths"
	"net/url"
)

func BackupHandler(_ *url.URL) error {
	return Backup()
}

func Backup() error {
	ea := nod.NewProgress("backing up metadata...")
	defer ea.End()

	amp, err := pathology.GetAbsDir(paths.Metadata)
	if err != nil {
		return ea.EndWithError(err)
	}

	abp, err := pathology.GetAbsDir(paths.Backups)
	if err != nil {
		return ea.EndWithError(err)
	}

	if err := packer.Pack(amp, abp, ea); err != nil {
		return ea.EndWithError(err)
	}

	ea.EndWithResult("done")

	return nil
}
