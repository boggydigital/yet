package rest

import (
	"net/http"
	"path"

	"github.com/boggydigital/strom"
	"github.com/boggydigital/strom/styles"
	"github.com/boggydigital/strom/vars/atoms"
	"github.com/boggydigital/strom/vars/calc"
	"github.com/boggydigital/strom/vars/colors"
	"github.com/boggydigital/strom/vars/sizes"
	"github.com/boggydigital/yet/data"
)

func GetChannelPlaylists(w http.ResponseWriter, r *http.Request) {

	// GET /channel_playlists/{channelId}

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

	// check if the channel has no playlists and refresh automatically
	if playlists, ok := rdx.GetAllValues(data.ChannelPlaylistsProperty, channelId); !ok || len(playlists) == 0 {
		http.Redirect(w, r, path.Join("/refresh_channel_playlists", channelId), http.StatusPermanentRedirect)
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
				styles.Decl("max-width", calc.Mult(sizes.XXXLarge, 4))))
	}

	channelMgmtRow := strom.Create("ul", atoms.FlexRowWrap(sizes.Small)...).
		Append(navButton("RSS", "https://www.youtube.com/feeds/videos.xml?channel_id="+channelId)).
		Append(navButton("Videos", path.Join("/channel", channelId))).
		Append(navButton("Refresh", path.Join("/refresh_channel_videos", channelId))).
		Append(navButton("Manage", path.Join("/manage_channel", channelId)))

	body.Append(channelMgmtRow)

	channelPlaylists := strom.Create("ul", atoms.FlexRowWrap(sizes.Normal)...)
	body.Append(channelPlaylists)

	if playlistIds, ok := rdx.GetAllValues(data.ChannelPlaylistsProperty, channelId); ok && len(playlistIds) > 0 {
		pl := new(playlistsList{playlistIds: playlistIds, rdx: rdx})
		channelPlaylists.Append(strom.OnDemand(pl.getPlaylistTiles))
	}

	if err = strom.WriteResponse(w, root); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
