package cli

import (
	"net/url"
	"strings"
)

func GetPlaylistMetadataHandler(u *url.URL) error {
	q := u.Query()
	ids := strings.Split(q.Get("id"), ",")
	force := q.Has("force")
	return GetPlaylistMetadata(force, ids...)
}

func GetPlaylistMetadata(force bool, ids ...string) error {
	return nil
}
