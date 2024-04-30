package yeti

import (
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yt_urls"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	playerUrlPfx = "/s/player/"
	playerUrlSfx = "/player_ias.vflset/en_US/base.js"
)

func PlayerVersion(playerUrl string) string {
	return strings.TrimPrefix(strings.TrimSuffix(playerUrl, playerUrlSfx), playerUrlPfx)
}

func GetPlayerContent(hc *http.Client, playerUrl string) (io.ReadCloser, error) {

	version := PlayerVersion(playerUrl)

	absPlayerPath, err := paths.AbsPlayerPath(version)
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(absPlayerPath); err == nil {
		return os.Open(absPlayerPath)
	}

	// local player doesn't exist - download and cache it

	pu := yt_urls.PlayerUrl(playerUrl)

	resp, err := hc.Get(pu.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	cachedPlayerFile, err := os.Create(absPlayerPath)
	if err != nil {
		return nil, err
	}
	defer cachedPlayerFile.Close()

	if _, err := io.Copy(cachedPlayerFile, resp.Body); err != nil {
		return nil, err
	}

	return os.Open(absPlayerPath)
}
