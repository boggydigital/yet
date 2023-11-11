package cli

import (
	"net/url"
	"strings"
)

func GetPosterHandler(u *url.URL) error {

	ids := strings.Split(u.Query().Get("id"), ",")
	return Download(ids)
}

func GetPoster(ids []string) error {
	return nil
}
