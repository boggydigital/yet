package rest

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
)

var (
	eprdx kvas.WriteableRedux
)

func Init() error {

	metadataDir, err := paths.GetAbsDir(paths.Metadata)
	if err != nil {
		return err
	}

	if eprdx, err = kvas.ReduxWriter(metadataDir, data.VideoEndedProperty, data.VideoProgressProperty); err != nil {
		return err
	}

	return err
}
