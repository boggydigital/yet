package rest

import (
	"github.com/boggydigital/redux"
	"github.com/boggydigital/yet/data"
	"maps"
	"net/http"
	"slices"
)

func GetUpdatePlaylist(w http.ResponseWriter, r *http.Request) {

	// GET /update_playlist?list

	var err error
	rdx, err = rdx.RefreshWriter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	q := r.URL.Query()

	playlistId := q.Get("list")

	if playlistId == "" {
		http.Redirect(w, r, "/list", http.StatusPermanentRedirect)
		return
	}

	boolPropertyInputs := map[string]string{
		data.PlaylistAutoRefreshProperty:  "auto-refresh",
		data.PlaylistExpandProperty:       "expand",
		data.PlaylistAutoDownloadProperty: "auto-download",
	}

	specialProperties := map[string]string{
		data.PlaylistDownloadPolicyProperty: "download-policy",
	}

	properties := slices.Collect(maps.Keys(boolPropertyInputs))
	properties = append(properties, slices.Collect(maps.Keys(specialProperties))...)

	for property, input := range boolPropertyInputs {
		if err := toggleProperty(playlistId, property, q.Has(input), rdx); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	for property, input := range specialProperties {
		switch property {
		case data.PlaylistDownloadPolicyProperty:
			policy := data.DefaultDownloadPolicy
			if dp := q.Get(input); dp != "" {
				policy = data.ParseDownloadPolicy(dp)
			}
			if err := rdx.ReplaceValues(data.PlaylistDownloadPolicyProperty, playlistId, string(policy)); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
	}

	http.Redirect(w, r, "/playlist?list="+playlistId, http.StatusTemporaryRedirect)

}

func toggleProperty(id, property string, condition bool, rdx redux.Writeable) error {
	if condition {
		if !rdx.HasValue(property, id, data.TrueValue) {
			return rdx.ReplaceValues(property, id, data.TrueValue)
		}
	} else {
		if rdx.HasKey(property, id) {
			return rdx.CutKeys(property, id)
		}
	}
	return nil
}
