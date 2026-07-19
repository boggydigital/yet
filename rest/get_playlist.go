package rest

import (
	"iter"
	"math"
	"net/http"
	"path"

	"github.com/boggydigital/redux"
	"github.com/boggydigital/strom"
	"github.com/boggydigital/strom/vars/atoms"
	"github.com/boggydigital/strom/vars/sizes"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
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
		http.Redirect(w, r, path.Join("/refresh_playlist", playlistId), http.StatusPermanentRedirect)
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
		Append(navButton("Refresh", path.Join("/refresh_playlist", playlistId))).
		Append(navButton("Manage", path.Join("/manage_playlist", playlistId)))

	body.Append(playlistMgmtRow)

	nepv := new(newEndedPlaylistVideos{playlistId: playlistId, rdx: rdx})

	body.Append(strom.OnDemand(nepv.getNewVideos))
	body.Append(strom.OnDemand(nepv.getEndedVideos))

	if err = strom.WriteResponse(w, root); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type newEndedPlaylistVideos struct {
	playlistId string
	rdx        redux.Readable
}

func (nepv *newEndedPlaylistVideos) getNewVideos() iter.Seq[strom.Element] {
	return nepv.getVideos(false)
}

func (nepv *newEndedPlaylistVideos) getEndedVideos() iter.Seq[strom.Element] {
	return nepv.getVideos(true)
}

func (nepv *newEndedPlaylistVideos) getVideos(ended bool) iter.Seq[strom.Element] {
	return func(yield func(element strom.Element) bool) {

		playlistVideos := strom.Create("ul", atoms.FlexRowWrap(sizes.Normal)...)
		if !ended {
			if newVideos := yeti.PlaylistNotEndedVideos(nepv.playlistId, math.MaxInt, nepv.rdx); len(newVideos) == 0 {
				return
			}
		}

		if plvs, ok := nepv.rdx.GetAllValues(data.PlaylistVideosProperty, nepv.playlistId); ok && len(plvs) > 0 {
			nev := new(newEndedVideos{ended: ended, videoIds: plvs, rdx: rdx})
			if ended {
				if !yield(strom.CreateText("h2", "Ended videos")) {
					return
				}
			}
			if !yield(playlistVideos.Append(strom.OnDemand(nev.getVideos))) {
				return
			}
		}
	}
}

type newEndedVideos struct {
	ended    bool
	videoIds []string
	rdx      redux.Readable
}

func (nev *newEndedVideos) getVideos() iter.Seq[strom.Element] {
	return func(yield func(element strom.Element) bool) {
		for _, videoId := range nev.videoIds {
			if nev.ended && !rdx.HasKey(data.VideoEndedDateProperty, videoId) {
				continue
			}
			if !nev.ended && rdx.HasKey(data.VideoEndedDateProperty, videoId) {
				continue
			}
			if !yield(videoTile(videoId, nev.rdx)) {
				return
			}
		}
	}
}
