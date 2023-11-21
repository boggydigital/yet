package cli

import (
	"net/url"
	"strings"
)

func GetMetadataHandler(u *url.URL) error {
	ids := strings.Split(u.Query().Get("id"), ",")
	return GetMetadata(ids...)
}

func GetMetadata(ids ...string) error {
	return nil
}
