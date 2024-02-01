package cli

import (
	"github.com/boggydigital/nod"
	"github.com/boggydigital/packer"
	"github.com/boggydigital/pasu"
	"github.com/boggydigital/yet/paths"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
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

	if err := packer.Pack(amp, abp, ea); err != nil {
		return ea.EndWithError(err)
	}

	ea.EndWithResult("done")

	if err := deleteOldBackups(abp); err != nil {
		return ea.EndWithError(err)
	}

	return nil
}

func deleteOldBackups(dir string) error {

	ofa := nod.Begin("deleting old backups...")
	defer ofa.End()

	d, err := os.Open(dir)
	if err != nil {
		return ofa.EndWithError(err)
	}
	defer d.Close()

	filenames, err := d.Readdirnames(-1)
	if err != nil {
		return ofa.EndWithError(err)
	}

	earliest := time.Now().Add(-daysToPreserveFiles * 24 * time.Hour)
	oldFiles := make([]string, 0)

	for _, fn := range filenames {

		fnse := fn
		for filepath.Ext(fnse) != "" {
			fnse = strings.TrimSuffix(fnse, filepath.Ext(fnse))
		}
		ft, err := time.Parse(nod.TimeFormat, fnse)
		if err != nil {
			nod.Log(err.Error())
			continue
		}

		if ft.After(earliest) {
			continue
		}

		oldFiles = append(oldFiles, fn)
	}

	if len(oldFiles) == 0 {
		ofa.EndWithResult("none found")
	} else {
		ofa.EndWithResult("found %d old backups", len(oldFiles))

		rofa := nod.NewProgress("removing old backups...")
		rofa.TotalInt(len(oldFiles))
		for _, fn := range oldFiles {
			filename := filepath.Join(dir, fn)
			if err := os.Remove(filename); err != nil {
				return rofa.EndWithError(err)
			}
			rofa.Increment()
		}
		rofa.EndWithResult("done")
	}

	return nil
}
