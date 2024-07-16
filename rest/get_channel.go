package rest

import (
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/rest/view_models"
	"net/http"
)

func GetChannel(w http.ResponseWriter, r *http.Request) {

	// GET /channel?channelId

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

	// check if the channel has no videos and refresh automatically
	if videos, ok := rdx.GetAllValues(data.ChannelVideosProperty, channelId); !ok || len(videos) == 0 {
		url := "/refresh_channel_videos?id=" + channelId
		http.Redirect(w, r, url, http.StatusPermanentRedirect)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	if err := tmpl.ExecuteTemplate(w, "channels", view_models.GetChannelViewModel(channelId, rdx)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
