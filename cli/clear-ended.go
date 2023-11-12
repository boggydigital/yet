package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"net/url"
	"strings"
)

func ClearEndedHandler(u *url.URL) error {
	ids := strings.Split(u.Query().Get("id"), ",")
	return ClearEnded(ids)
}

func ClearEnded(ids []string) error {
	cea := nod.Begin("clearing ended...")
	defer cea.End()

	metadataDir, err := paths.GetAbsDir(paths.Metadata)
	if err != nil {
		return cea.EndWithError(err)
	}

	rdx, err := kvas.ConnectRedux(metadataDir, data.VideoEndedProperty)
	if err != nil {
		return cea.EndWithError(err)
	}

	endedClear := make(map[string][]string)
	for _, id := range ids {
		endedClear[id] = nil
	}

	if err := rdx.BatchReplaceValues(endedClear); err != nil {
		return cea.EndWithError(err)
	}

	cea.EndWithResult("done")

	return nil
}
