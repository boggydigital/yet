package rest

import (
	"net/http"
	"path"

	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
)

func GetRefreshChannelVideos(w http.ResponseWriter, r *http.Request) {

	// GET /refresh_channel_videos/{channelId}

	var err error
	rdx, err = rdx.RefreshWriter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	channelId := r.PathValue("channelId")

	if channelId == "" {
		http.Redirect(w, r, "/list", http.StatusPermanentRedirect)
		return
	}

	expand := false
	if exp, ok := rdx.GetLastVal(data.ChannelExpandProperty, channelId); ok && exp == data.TrueValue {
		expand = true
	}

	if err = yeti.GetChannelVideosMetadata(nil, channelId, expand, rdx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, path.Join("/channel", channelId), http.StatusTemporaryRedirect)
}
