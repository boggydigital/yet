package rest

import (
	"iter"
	"net/http"
	"path"

	"github.com/boggydigital/redux"
	"github.com/boggydigital/strom"
	"github.com/boggydigital/strom/vars/atoms"
	"github.com/boggydigital/strom/vars/calc"
	"github.com/boggydigital/strom/vars/colors"
	"github.com/boggydigital/strom/vars/sizes"
	"github.com/boggydigital/yet/data"
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
			SetStyle(map[string]string{
				"color":     colors.Gray,
				"max-width": calc.Mult(sizes.XXXLarge, 4),
			}))
	}

	channelNavButtonsRow := strom.Create("ul", atoms.FlexRowWrap(sizes.Small)...).
		Append(navButton("RSS", "https://www.youtube.com/feeds/videos.xml?channel_id="+channelId)).
		Append(navButton("Playlists", path.Join("/channel_playlists", channelId))).
		Append(navButton("Refresh", path.Join("/refresh_channel_videos", channelId))).
		Append(navButton("Manage", path.Join("/manage_channel", channelId)))

	body.Append(channelNavButtonsRow)

	cv := new(channelVideos{channelId: channelId, rdx: rdx})

	videos := strom.Create("ul", atoms.FlexRowWrap(sizes.Normal)...)
	body.Append(videos)

	videos.Append(strom.OnDemand(cv.getVideos))

	if err = strom.WriteResponse(w, root); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

type channelVideos struct {
	channelId string
	rdx       redux.Readable
}

func (cv *channelVideos) getVideos() iter.Seq[strom.Element] {
	return func(yield func(element strom.Element) bool) {

		if chvs, ok := rdx.GetAllValues(data.ChannelVideosProperty, cv.channelId); ok && len(chvs) > 0 {
			for _, videoId := range chvs {
				if !yield(videoTile(videoId, rdx)) {
					return
				}
			}
		}
	}
}
