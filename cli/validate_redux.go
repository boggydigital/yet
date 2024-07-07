package cli

import (
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/pathways"
	"github.com/boggydigital/yet/paths"
)

func validateWritableRedux(rdx kevlar.WriteableRedux, properties ...string) (kevlar.WriteableRedux, error) {
	if rdx == nil {
		metadataDir, err := pathways.GetAbsDir(paths.Metadata)
		if err != nil {
			return nil, err
		}

		rdx, err = kevlar.NewReduxWriter(metadataDir, properties...)
		if err != nil {
			return nil, err
		}
	}
	if err := rdx.MustHave(properties...); err != nil {
		return nil, err
	}

	return rdx, nil
}
