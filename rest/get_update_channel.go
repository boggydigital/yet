package rest

import (
	"maps"
	"net/http"
	"path"
	"slices"

	"github.com/boggydigital/yet/data"
)

func GetUpdateChannel(w http.ResponseWriter, r *http.Request) {

	// GET /update_channel/{channelId}

	var err error
	rdx, err = rdx.RefreshWriter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	q := r.URL.Query()

	channelId := r.PathValue("channelId")

	if channelId == "" {
		http.Redirect(w, r, "/list", http.StatusPermanentRedirect)
		return
	}

	boolPropertyInputs := map[string]string{
		data.ChannelAutoRefreshProperty:  "auto-refresh",
		data.ChannelExpandProperty:       "expand",
		data.ChannelAutoDownloadProperty: "auto-download",
	}

	specialProperties := map[string]string{
		data.ChannelDownloadPolicyProperty: "download-policy",
	}

	properties := slices.Collect(maps.Keys(boolPropertyInputs))
	properties = append(properties, slices.Collect(maps.Keys(specialProperties))...)

	for property, input := range boolPropertyInputs {
		if err = toggleProperty(channelId, property, q.Has(input), rdx); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	for property, input := range specialProperties {
		switch property {
		case data.ChannelDownloadPolicyProperty:
			policy := data.DefaultDownloadPolicy
			if dp := q.Get(input); dp != "" {
				policy = data.ParseDownloadPolicy(dp)
			}
			if err = rdx.ReplaceValues(data.ChannelDownloadPolicyProperty, channelId, string(policy)); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
	}

	http.Redirect(w, r, path.Join("/channel", channelId), http.StatusTemporaryRedirect)
}
