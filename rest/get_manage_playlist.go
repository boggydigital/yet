package rest

import (
	"net/http"
	"path"

	"github.com/boggydigital/strom"
	"github.com/boggydigital/strom/styles"
	"github.com/boggydigital/strom/vars/atoms"
	"github.com/boggydigital/strom/vars/colors"
	"github.com/boggydigital/strom/vars/sizes"
	"github.com/boggydigital/yet/data"
)

func GetManagePlaylist(w http.ResponseWriter, r *http.Request) {

	// GET /manage_playlist?list

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

	var playlistTitle string
	if pt, ok := rdx.GetLastVal(data.PlaylistTitleProperty, playlistId); ok && pt != "" {
		playlistTitle = pt
	}

	root, body := strom.RootBody(playlistTitle, atoms.FlexCol(sizes.Normal)...)

	topRow := strom.Create("ul", atoms.FlexRow(sizes.Small)...).AddAtom(atoms.AlignItemsCenter)
	body.Append(topRow)

	topRow.Append(navButton("Home", "/"))
	topRow.Append(strom.CreateText("h2", "Manage playlist"))

	body.Append(playlistTile(playlistId, rdx))

	originRow := strom.Create("ul", atoms.FlexRow(sizes.Small)...).
		AddAtom(atoms.AlignItemsCenter)
	body.Append(originRow)

	originRow.Append(
		navButton("Origin", "https://www.youtube.com/playlist?list="+playlistId),
		navButton("RSS", "https://www.youtube.com/feeds/videos.xml?playlist_id="+playlistId),
		strom.CreateText("span", "Playlist ID").
			SetStyle(styles.Decl("color", colors.Gray)),
		strom.CreateText("span", playlistId, atoms.FontWeightBold))

	form := strom.Create("form", atoms.FlexColWrap(sizes.Normal)...).
		SetAttribute("id", "manage-playlist").
		SetAttribute("method", "get").
		SetAttribute("action", path.Join("/update_playlist/", playlistId))
	body.Append(form)

	autoRefresh := rdx.HasKey(data.PlaylistAutoRefreshProperty, playlistId)
	form.Append(switchTitleSubtitle(autoRefresh, "auto-refresh", "Auto refresh", "Update metadata, videos."))

	expandPlaylist := rdx.HasKey(data.PlaylistExpandProperty, playlistId)
	form.Append(switchTitleSubtitle(expandPlaylist, "expand", "Expand channel videos", "On: Get all videos in a playlist. Off: Only get the latest 100 videos."))

	autoDownload := rdx.HasKey(data.PlaylistAutoDownloadProperty, playlistId)
	form.Append(switchTitleSubtitle(autoDownload, "auto-download", "Auto download videos", "Download new videos."))

	var downloadPolicy data.DownloadPolicy
	if dps, ok := rdx.GetLastVal(data.PlaylistDownloadPolicyProperty, playlistId); ok && dps != "" {
		downloadPolicy = data.ParseDownloadPolicy(dps)
	}
	form.Append(downloadPolicySelect(downloadPolicy))

	body.Append(submitButton("Update", "manage-playlist"))

	if err = strom.WriteResponse(w, root); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
