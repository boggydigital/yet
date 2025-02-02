package cli

import (
	"github.com/boggydigital/pathways"
	"github.com/boggydigital/redux"
	"github.com/boggydigital/yet/data"
)

func validateWritableRedux(rdx redux.Writeable, properties ...string) (redux.Writeable, error) {
	if rdx == nil {
		metadataDir, err := pathways.GetAbsDir(data.Metadata)
		if err != nil {
			return nil, err
		}

		rdx, err = redux.NewWriter(metadataDir, properties...)
		if err != nil {
			return nil, err
		}
	}
	if err := rdx.MustHave(properties...); err != nil {
		return nil, err
	}

	var err error
	if rdx, err = rdx.RefreshWriter(); err != nil {
		return nil, err
	}

	return rdx, nil
}
