package rest

import (
	"iter"
	"math"
	"net/http"
	"path"

	"github.com/boggydigital/redux"
	"github.com/boggydigital/strom"
	"github.com/boggydigital/strom/vars/atoms"
	"github.com/boggydigital/strom/vars/colors"
	"github.com/boggydigital/strom/vars/sizes"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
)

func GetChannel(w http.ResponseWriter, r *http.Request) {

	// GET /channel/{channelId}

	var err error
	rdx, err = rdx.RefreshWriter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	channelId := r.PathValue("channelId")

	if channelId == "" {
		http.Redirect(w, r, "/list", http.StatusPermanentRedirect)
		return
	}

	var channelTitle string
	if ct, ok := rdx.GetLastVal(data.ChannelTitleProperty, channelId); ok && ct != "" {
		channelTitle = ct
	}

	root, body := strom.RootBody(channelTitle, atoms.FlexCol(sizes.Normal)...)

	topRow := strom.Create("ul", atoms.FlexRow(sizes.Small)...).AddAtom(atoms.AlignItemsCenter)
	body.Append(topRow)

	topRow.Append(navButton("Home", "/"))
	topRow.Append(strom.CreateText("h2", "Channel"))

	body.Append(channelTile(channelId, rdx))

	refreshChannelForm := strom.Create("form").
		SetAttribute("id", "refresh-channel").
		SetAttribute("method", "post").
		SetAttribute("action", path.Join("/refresh_channel/", channelId)).
		SetStyle("display:none")
	body.Append(refreshChannelForm)

	channelMgmtRow := strom.Create("ul", atoms.FlexRowWrap(sizes.Small)...).Append(
		navButton("Manage", path.Join("/manage_channel", channelId), colors.Blue),
		navButton("Playlists", "#channel_playlists"),
		submitButton("Refresh", "refresh-channel", colors.Green))

	body.Append(channelMgmtRow)

	cv := new(newEndedChannelVideos{channelId: channelId, rdx: rdx})

	body.Append(strom.OnDemand(cv.getNewVideos))
	body.Append(strom.OnDemand(cv.getEndedVideos))

	body.Append(highVisibilityAnchor("Channel playlists"))

	if playlistIds, ok := rdx.GetAllValues(data.ChannelPlaylistsProperty, channelId); ok && len(playlistIds) > 0 {
		channelPlaylists := strom.Create("ul", atoms.FlexRowWrap(sizes.Normal)...)
		body.Append(channelPlaylists)

		pl := new(playlistsList{playlistIds: playlistIds, rdx: rdx})
		channelPlaylists.Append(strom.OnDemand(pl.getPlaylistTiles))
	} else {
		body.Append(strom.CreateText("span", "Channel has no playlists").
			SetStyle("color:" + colors.Gray))
	}

	if err = strom.WriteResponse(w, root); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

type newEndedChannelVideos struct {
	channelId string
	rdx       redux.Readable
}

func (necv *newEndedChannelVideos) getNewVideos() iter.Seq[strom.Element] {
	return necv.getVideos(false)
}

func (necv *newEndedChannelVideos) getEndedVideos() iter.Seq[strom.Element] {
	return necv.getVideos(true)
}

func (necv *newEndedChannelVideos) getVideos(ended bool) iter.Seq[strom.Element] {
	return func(yield func(element strom.Element) bool) {

		channelVideos := strom.Create("ul", atoms.FlexRowWrap(sizes.Normal)...)
		if !ended {
			if newVideos := yeti.ChannelNotEndedVideos(necv.channelId, math.MaxInt, necv.rdx); len(newVideos) == 0 {
				return
			}
		}

		if chvs, ok := rdx.GetAllValues(data.ChannelVideosProperty, necv.channelId); ok && len(chvs) > 0 {
			nev := new(newEndedVideos{ended: ended, videoIds: chvs, rdx: rdx})

			if ended {
				if !yield(highVisibilityAnchor("Ended videos")) {
					return
				}
			}

			if !yield(channelVideos.Append(strom.OnDemand(nev.getVideos))) {
				return
			}
		}
	}
}
