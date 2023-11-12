package rest

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
)

var progressRdx kvas.ReduxValues

func Init() error {

	metadataDir, err := paths.GetAbsDir(paths.Metadata)
	if err != nil {
		return err
	}

	progressRdx, err = kvas.ConnectRedux(metadataDir, data.VideoProgressProperty)

	return err
}
