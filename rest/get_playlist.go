package rest

import (
	"iter"
	"net/http"
	"path"

	"github.com/boggydigital/redux"
	"github.com/boggydigital/strom"
	"github.com/boggydigital/strom/vars/atoms"
	"github.com/boggydigital/strom/vars/sizes"
	"github.com/boggydigital/yet/data"
)

func GetPlaylist(w http.ResponseWriter, r *http.Request) {

	// GET /playlist/{playlistId}

	var err error
	rdx, err = rdx.RefreshWriter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	playlistId := r.PathValue("playlistId")

	if playlistId == "" {
		http.Redirect(w, r, "/list", http.StatusPermanentRedirect)
		return
	}

	// check if the playlist has no videos and refresh automatically
	if videos, ok := rdx.GetAllValues(data.PlaylistVideosProperty, playlistId); !ok || len(videos) == 0 {
		url := "/refresh_playlist?list=" + playlistId
		http.Redirect(w, r, url, http.StatusPermanentRedirect)
		return
	}

	var playlistTitle string
	if pt, ok := rdx.GetLastVal(data.PlaylistTitleProperty, playlistId); ok && pt != "" {
		playlistTitle = pt
	}

	root, body := strom.RootBody(playlistTitle, atoms.FlexCol(sizes.Normal)...)

	topRow := strom.Create("ul", atoms.FlexRow(sizes.Small)...).AddAtom(atoms.AlignItemsCenter)
	body.Append(topRow)

	topRow.Append(navButton("Home", "/"))
	topRow.Append(strom.CreateText("h2", "Playlist"))

	body.Append(playlistTile(playlistId, rdx))

	playlistMgmtRow := strom.Create("ul", atoms.FlexRowWrap(sizes.Small)...).
		Append(navButton("RSS", "https://www.youtube.com/feeds/videos.xml?playlist_id="+playlistId)).
		Append(navButton("Refresh", path.Join("/refresh_playlist?list="+playlistId))).
		Append(navButton("Manage", path.Join("/manage_playlist?list=", playlistId)))

	body.Append(playlistMgmtRow)

	pv := new(playlistVideos{playlistId: playlistId, rdx: rdx})

	videos := strom.Create("ul", atoms.FlexRowWrap(sizes.Normal)...)
	body.Append(videos)

	videos.Append(strom.OnDemand(pv.getVideos))

	if err = strom.WriteResponse(w, root); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	//w.Header().Set("Content-Type", "text/html")
	//
	//if err := tmpl.ExecuteTemplate(w, "playlists", view_models.GetPlaylistViewModel(playlistId, rdx)); err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
}

type playlistVideos struct {
	playlistId string
	rdx        redux.Readable
}

func (pv *playlistVideos) getVideos() iter.Seq[strom.Element] {
	return func(yield func(element strom.Element) bool) {
		if plvs, ok := rdx.GetAllValues(data.PlaylistVideosProperty, pv.playlistId); ok && len(plvs) > 0 {
			for _, videoId := range plvs {
				if !yield(videoTile(videoId, rdx)) {
					return
				}
			}
		}
	}
}
