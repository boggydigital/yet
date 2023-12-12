package cli

import (
	"github.com/boggydigital/yet/data"
	"net/url"
	"strings"
)

func AddPlaylistsHandler(u *url.URL) error {
	q := u.Query()

	watchlist := strings.Split(q.Get("watchlist"), ",")

	return AddPlaylists(map[string][]string{
		data.VideosWatchlistProperty: watchlist,
	})
}

func AddPlaylists(propertyValues map[string][]string) error {
	return nil
}
