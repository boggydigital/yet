package rest

import (
	"maps"
	"net/http"
	"path"
	"slices"

	"github.com/boggydigital/redux"
	"github.com/boggydigital/yet/data"
)

func PostUpdatePlaylist(w http.ResponseWriter, r *http.Request) {

	// POST /update_playlist/{playlistId}

	var err error
	rdx, err = rdx.RefreshWriter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	playlistId := r.PathValue("playlistId")

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
		if err = toggleProperty(playlistId, property, r.Form.Get(input) == "on", rdx); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	for property, input := range specialProperties {
		switch property {
		case data.PlaylistDownloadPolicyProperty:
			policy := data.DefaultDownloadPolicy
			if dp := r.Form.Get(input); dp != "" {
				policy = data.ParseDownloadPolicy(dp)
			}
			if err = rdx.ReplaceValues(data.PlaylistDownloadPolicyProperty, playlistId, string(policy)); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
	}

	w.Header().Set("Location", path.Join("/playlist", playlistId))
	w.WriteHeader(http.StatusSeeOther)
	return
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
