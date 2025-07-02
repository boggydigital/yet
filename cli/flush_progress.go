package cli

import (
	"github.com/boggydigital/redux"
	"github.com/boggydigital/yet/data"
)

func flushProgress(rdx redux.Writeable) error {

	var err error
	if rdx, err = validateWritableRedux(rdx, data.VideoProperties()...); err != nil {
		return err
	}

	data.ProgressMux.Lock()
	defer data.ProgressMux.Unlock()

	if err = rdx.BatchReplaceValues(data.VideoProgressProperty, data.VideosProgress); err != nil {
		return err
	}

	data.VideosProgress = make(map[string][]string)

	return nil
}
