package rest

import (
	"github.com/boggydigital/camino"
	"github.com/boggydigital/redux"
	"github.com/boggydigital/yet/data"
)

var (
	rdx redux.Writeable
)

func Init() error {

	metadataDir := camino.GetAbs(data.Metadata)

	var err error
	if rdx, err = redux.NewWriter(metadataDir, data.AllProperties()...); err != nil {
		return err
	}

	return nil
}
