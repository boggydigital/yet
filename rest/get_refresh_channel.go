package rest

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yet/yeti"
	"net/http"
)

func GetRefreshChannel(w http.ResponseWriter, r *http.Request) {

	// GET /refresh_channel?id

	channelId := r.URL.Query().Get("id")

	if channelId == "" {
		http.Redirect(w, r, "/new", http.StatusPermanentRedirect)
		return
	}

	metadataDir, err := paths.GetAbsDir(paths.Metadata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	chRdx, err := kvas.NewReduxWriter(metadataDir, data.AllProperties()...)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := yeti.GetChannelPageMetadata(nil, channelId, chRdx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/channel?id="+channelId, http.StatusTemporaryRedirect)
}
