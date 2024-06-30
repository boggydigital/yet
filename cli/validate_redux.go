package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/pathways"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
)

func validateWritableRedux(rdx kvas.WriteableRedux, properties ...string) (kvas.WriteableRedux, error) {
	if rdx == nil {
		metadataDir, err := pathways.GetAbsDir(paths.Metadata)
		if err != nil {
			return nil, err
		}

		rdx, err = kvas.NewReduxWriter(metadataDir, properties...)
		if err != nil {
			return nil, err
		}
	} else if err := rdx.MustHave(data.PlaylistProperties()...); err != nil {
		return nil, err
	}

	return rdx, nil
}
