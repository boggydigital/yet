package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func SetEndedHandler(u *url.URL) error {
	ids := strings.Split(u.Query().Get("id"), ",")
	return SetEnded(ids)
}

func SetEnded(ids []string) error {

	sea := nod.Begin("setting ended...")
	defer sea.End()

	metadataDir, err := paths.GetAbsDir(paths.Metadata)
	if err != nil {
		return sea.EndWithError(err)
	}

	rdx, err := kvas.ConnectRedux(metadataDir, data.VideoEndedProperty)
	if err != nil {
		return sea.EndWithError(err)
	}

	endedSet := make(map[string][]string)
	for _, id := range ids {
		endedSet[id] = []string{time.Now().Format(http.TimeFormat)}
	}

	if err := rdx.BatchReplaceValues(endedSet); err != nil {
		return sea.EndWithError(err)
	}

	sea.EndWithResult("done")

	return nil
}
