package yeti

import (
	"errors"
	"github.com/boggydigital/yt_urls"
	"net/http"
)

func DecodeSignatureCipher(hc *http.Client, ipr *yt_urls.InitialPlayerResponse) error {
	if !ipr.SignatureCipher() {
		return nil
	}
	return errors.New("signatureCipher not implemented")
}
