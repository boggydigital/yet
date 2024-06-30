package rest

import (
	"fmt"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/pathways"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"net/http"
	"net/url"
)

func GetUpdatePlaylist(w http.ResponseWriter, r *http.Request) {

	// GET /update_playlist?list

	playlistId := r.URL.Query().Get("list")

	if playlistId == "" {
		http.Redirect(w, r, "/list", http.StatusPermanentRedirect)
		return
	}

	metadataDir, err := pathways.GetAbsDir(paths.Metadata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	properties := []string{
		data.PlaylistAutoRefreshProperty,
		data.PlaylistAutoDownloadProperty,
		data.PlaylistPreferSingleFormatProperty,
	}

	plRdx, err := kvas.NewReduxWriter(metadataDir, properties...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, p := range properties {
		if err := updatePlaylistProperty(playlistId, p, r.URL, plRdx); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	http.Redirect(w, r, "/playlist?list="+playlistId, http.StatusTemporaryRedirect)

}

func updatePlaylistProperty(playlistId string, property string, u *url.URL, rdx kvas.WriteableRedux) error {

	flagStr := ""
	switch property {
	case data.PlaylistAutoRefreshProperty:
		flagStr = "auto-refresh"
	case data.PlaylistAutoDownloadProperty:
		flagStr = "auto-download"
	case data.PlaylistPreferSingleFormatProperty:
		flagStr = "prefer-single-format"
	default:
		return fmt.Errorf("unsupported property %s", property)
	}

	flag := u.Query().Has(flagStr)

	var err error

	if flag {
		if !rdx.HasKey(property, playlistId) {
			err = rdx.AddValues(property, playlistId, data.TrueValue)
		}
	} else {
		if rdx.HasKey(property, playlistId) {
			err = rdx.CutKeys(property, playlistId)
		}
	}

	return err
}
