package rest

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/pathways"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"golang.org/x/exp/maps"
	"net/http"
)

func GetUpdatePlaylist(w http.ResponseWriter, r *http.Request) {

	// GET /update_playlist?list

	q := r.URL.Query()

	playlistId := q.Get("list")

	if playlistId == "" {
		http.Redirect(w, r, "/list", http.StatusPermanentRedirect)
		return
	}

	metadataDir, err := pathways.GetAbsDir(paths.Metadata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	propertyInputs := map[string]string{
		data.PlaylistAutoRefreshProperty:        "auto-refresh",
		data.PlaylistExpandProperty:             "expand",
		data.PlaylistAutoDownloadProperty:       "auto-download",
		data.PlaylistPreferSingleFormatProperty: "prefer-single-format",
	}

	properties := maps.Keys(propertyInputs)
	properties = append(properties, data.PlaylistDownloadPolicyProperty)

	plRdx, err := kvas.NewReduxWriter(metadataDir, properties...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for property, input := range propertyInputs {
		if err := toggleProperty(playlistId, property, q.Has(input), plRdx); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	// download policy requires special non-binary handling
	policy := data.Unset
	if q.Has("download-policy-all") {
		policy = data.All
	} else {
		policy = data.Unset
	}

	if err := plRdx.ReplaceValues(data.PlaylistDownloadPolicyProperty, playlistId, string(policy)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/playlist?list="+playlistId, http.StatusTemporaryRedirect)

}

func toggleProperty(playlistId, property string, condition bool, rdx kvas.WriteableRedux) error {
	if condition {
		return rdx.ReplaceValues(property, playlistId, data.TrueValue)
	} else {
		return rdx.CutKeys(property, playlistId)
	}
}
