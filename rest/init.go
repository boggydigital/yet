package rest

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
)

var (
	endedRdx    kvas.ReduxValues
	progressRdx kvas.ReduxValues
)

func Init() error {

	metadataDir, err := paths.GetAbsDir(paths.Metadata)
	if err != nil {
		return err
	}

	if progressRdx, err = kvas.ConnectRedux(metadataDir, data.VideoProgressProperty); err != nil {
		return err
	}

	if endedRdx, err = kvas.ConnectRedux(metadataDir, data.VideoEndedProperty); err != nil {
		return err
	}

	return err
}
