package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"net/url"
	"strings"
)

func ClearProgressHandler(u *url.URL) error {
	ids := strings.Split(u.Query().Get("id"), ",")
	return ClearProgress(ids)
}

func ClearProgress(ids []string) error {
	cpa := nod.Begin("clearing progress...")
	defer cpa.End()

	metadataDir, err := paths.GetAbsDir(paths.Metadata)
	if err != nil {
		return cpa.EndWithError(err)
	}

	rdx, err := kvas.ConnectRedux(metadataDir, data.VideoProgressProperty)
	if err != nil {
		return cpa.EndWithError(err)
	}

	progressClear := make(map[string][]string)
	for _, id := range ids {
		progressClear[id] = nil
	}

	if err := rdx.BatchReplaceValues(progressClear); err != nil {
		return cpa.EndWithError(err)
	}

	cpa.EndWithResult("done")

	return nil
}
