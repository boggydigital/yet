package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pathology"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yet/yeti"
	"golang.org/x/exp/maps"
	"net/url"
	"strings"
)

func RemovePlaylistsHandler(u *url.URL) error {
	q := u.Query()

	watchlist := strings.Split(q.Get("watchlist"), ",")
	downloadQueue := strings.Split(q.Get("download-queue"), ",")

	return RemovePlaylists(map[string][]string{
		data.PlaylistWatchlistProperty:     watchlist,
		data.PlaylistDownloadQueueProperty: downloadQueue,
	})
}

func RemovePlaylists(propertyValues map[string][]string) error {
	rpa := nod.NewProgress("removing playlists...")
	defer rpa.End()

	metadataDir, err := pathology.GetAbsDir(paths.Metadata)
	if err != nil {
		return rpa.EndWithError(err)
	}

	rdx, err := kvas.NewReduxWriter(metadataDir, maps.Keys(propertyValues)...)
	if err != nil {
		return rpa.EndWithError(err)
	}

	rpa.TotalInt(len(propertyValues))

	for property, values := range propertyValues {
		if err := removePropertyValues(rdx, yeti.ParsePlaylistIds, property, values...); err != nil {
			return rpa.EndWithError(err)
		}
		rpa.Increment()
	}

	rpa.EndWithResult("done")

	return nil
}
