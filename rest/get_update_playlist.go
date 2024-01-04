package rest

import (
	"fmt"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"net/http"
	"net/url"
)

func GetUpdatePlaylist(w http.ResponseWriter, r *http.Request) {

	// GET /update_playlist?list

	var err error
	rdx, err = rdx.RefreshReader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	playlistId := r.URL.Query().Get("list")

	if playlistId == "" {
		http.Redirect(w, r, "/list", http.StatusPermanentRedirect)
		return
	}

	metadataDir, err := paths.GetAbsDir(paths.Metadata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	plRdx, err := kvas.NewReduxWriter(metadataDir, data.PlaylistWatchlistProperty, data.PlaylistDownloadQueueProperty)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := updatePlaylistProperty(playlistId, data.PlaylistWatchlistProperty, r.URL, plRdx); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := updatePlaylistProperty(playlistId, data.PlaylistDownloadQueueProperty, r.URL, plRdx); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/playlist?list="+playlistId, http.StatusTemporaryRedirect)

}

func updatePlaylistProperty(playlistId string, property string, u *url.URL, rdx kvas.WriteableRedux) error {

	flagStr := ""
	switch property {
	case data.PlaylistWatchlistProperty:
		flagStr = "refresh"
	case data.PlaylistDownloadQueueProperty:
		flagStr = "download"
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
			err = rdx.CutValues(property, playlistId, data.TrueValue)
		}
	}

	return err
}
