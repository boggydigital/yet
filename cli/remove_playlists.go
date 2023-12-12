package cli

import (
	"github.com/boggydigital/yet/data"
	"net/url"
	"strings"
)

func RemovePlaylistsHandler(u *url.URL) error {
	q := u.Query()

	watchlist := strings.Split(q.Get("watchlist"), ",")

	return RemovePlaylists(map[string][]string{
		data.VideosWatchlistProperty: watchlist,
	})
}

func RemovePlaylists(propertyValues map[string][]string) error {
	return nil
}
