package rest

import (
	"iter"
	"math"
	"net/http"
	"path"

	"github.com/boggydigital/redux"
	"github.com/boggydigital/strom"
	"github.com/boggydigital/strom/styles"
	"github.com/boggydigital/strom/vars/atoms"
	"github.com/boggydigital/strom/vars/calc"
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

	// check if the channel has no videos and refresh automatically
	if videos, ok := rdx.GetAllValues(data.ChannelVideosProperty, channelId); !ok || len(videos) == 0 {
		http.Redirect(w, r, path.Join("/refresh_channel_videos", channelId), http.StatusPermanentRedirect)
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

	if cd, ok := rdx.GetLastVal(data.ChannelDescriptionProperty, channelId); ok && cd != "" {
		body.Append(strom.CreateText("span", cd).
			SetStyle(
				styles.Decl("color", colors.Gray),
				styles.Decl("max-width", calc.Mult(sizes.XXXLarge, 4)),
				"word-break:break-word"))
	}

	channelMgmtRow := strom.Create("ul", atoms.FlexRowWrap(sizes.Small)...).
		Append(navButton("Refresh", path.Join("/refresh_channel_videos", channelId))).
		Append(navButton("Playlists", path.Join("/channel_playlists", channelId))).
		Append(navButton("Manage", path.Join("/manage_channel", channelId)))

	body.Append(channelMgmtRow)

	cv := new(newEndedChannelVideos{channelId: channelId, rdx: rdx})

	body.Append(strom.OnDemand(cv.getNewVideos))
	body.Append(strom.OnDemand(cv.getEndedVideos))

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
				if !yield(strom.CreateText("h2", "Ended videos")) {
					return
				}
			}

			if !yield(channelVideos.Append(strom.OnDemand(nev.getVideos))) {
				return
			}
		}
	}
}
