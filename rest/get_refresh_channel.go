package rest

import (
	"github.com/boggydigital/yet/yeti"
	"net/http"
)

func GetRefreshChannel(w http.ResponseWriter, r *http.Request) {

	// GET /refresh_channel?id

	var err error
	rdx, err = rdx.RefreshWriter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	channelId := r.URL.Query().Get("id")

	if channelId == "" {
		http.Redirect(w, r, "/list", http.StatusPermanentRedirect)
		return
	}

	if err := yeti.GetChannelPageMetadata(nil, channelId, rdx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/channel?id="+channelId, http.StatusTemporaryRedirect)
}
