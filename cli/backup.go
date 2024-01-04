package cli

import (
	"github.com/boggydigital/nod"
	"github.com/boggydigital/packer"
	"github.com/boggydigital/pathology"
	"github.com/boggydigital/yet/paths"
	"net/url"
)

func BackupHandler(u *url.URL) error {
	aofp, err := pathology.GetAbsDir(paths.Backups)
	if err != nil {
		return err
	}
	return Backup(aofp)
}

func Backup(to string) error {
	ea := nod.NewProgress("backing up metadata...")
	defer ea.End()

	amp, err := pathology.GetAbsDir(paths.Metadata)
	if err != nil {
		return ea.EndWithError(err)
	}

	if err := packer.Pack(amp, to, ea); err != nil {
		return ea.EndWithError(err)
	}

	ea.EndWithResult("done")

	return nil
}
