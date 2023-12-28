package rest

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
)

var (
	rdx kvas.ReadableRedux
)

func Init() error {

	metadataDir, err := paths.GetAbsDir(paths.Metadata)
	if err != nil {
		return err
	}

	if rdx, err = kvas.NewReduxReader(metadataDir, data.AllProperties()...); err != nil {
		return err
	}

	return err
}
