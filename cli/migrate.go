package cli

import (
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pathways"
	"github.com/boggydigital/yet/data"
	"net/url"
)

func MigrateHandler(_ *url.URL) error {
	return Migrate()
}

func Migrate() error {
	ma := nod.Begin("migrating data...")
	defer ma.Done()

	if err := Backup(); err != nil {
		return err
	}

	amd, err := pathways.GetAbsDir(data.Metadata)
	if err != nil {
		return err
	}

	if err := kevlar.Migrate(amd); err != nil {
		return err
	}

	return nil
}
