package rest

import (
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/yet/data"
	"golang.org/x/exp/maps"
	"net/http"
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
		data.PlaylistAutoRefreshProperty:        "auto-refresh",
		data.PlaylistExpandProperty:             "expand",
		data.PlaylistAutoDownloadProperty:       "auto-download",
		data.PlaylistPreferSingleFormatProperty: "prefer-single-format",
	}

	specialProperties := map[string]string{
		data.PlaylistDownloadPolicyProperty: "download-policy",
	}

	properties := maps.Keys(boolPropertyInputs)
	properties = append(properties, maps.Keys(specialProperties)...)

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
				policy = data.ParsePlaylistDownloadPolicy(dp)
			}
			if err := rdx.ReplaceValues(data.PlaylistDownloadPolicyProperty, playlistId, string(policy)); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
	}

	http.Redirect(w, r, "/playlist?list="+playlistId, http.StatusTemporaryRedirect)

}

func toggleProperty(id, property string, condition bool, rdx kevlar.WriteableRedux) error {
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
