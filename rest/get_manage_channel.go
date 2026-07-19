package rest

import (
	"net/http"
	"path"

	"github.com/boggydigital/strom"
	"github.com/boggydigital/strom/styles"
	"github.com/boggydigital/strom/vars/atoms"
	"github.com/boggydigital/strom/vars/colors"
	"github.com/boggydigital/strom/vars/font_sizes"
	"github.com/boggydigital/strom/vars/sizes"
	"github.com/boggydigital/yet/data"
)

func GetManageChannel(w http.ResponseWriter, r *http.Request) {

	// GET /manage_channel/{channelId}

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
	topRow.Append(strom.CreateText("h2", "Manage channel"))

	body.Append(channelTile(channelId, rdx))

	originRow := strom.Create("ul", atoms.FlexRowWrap(sizes.Small)...).
		AddAtom(atoms.AlignItemsCenter)
	body.Append(originRow)

	originRow.Append(
		navButton("Origin", path.Join("https://www.youtube.com/channel", channelId)),
		navButton("RSS", "https://www.youtube.com/feeds/videos.xml?channel_id="+channelId))

	form := strom.Create("form", atoms.FlexColWrap(sizes.Normal)...).
		SetAttribute("id", "manage-channel").
		SetAttribute("method", "get").
		SetAttribute("action", path.Join("/update_channel/", channelId))
	body.Append(form)

	autoRefresh := rdx.HasKey(data.ChannelAutoRefreshProperty, channelId)
	form.Append(switchTitleSubtitle(autoRefresh, "auto-refresh", "Auto refresh", "Update metadata, videos."))

	expandChannel := rdx.HasKey(data.ChannelExpandProperty, channelId)
	form.Append(switchTitleSubtitle(expandChannel, "expand", "Expand channel videos", "On: Get all videos in a channel. Off: Only get the latest 30 videos."))

	autoDownload := rdx.HasKey(data.ChannelAutoDownloadProperty, channelId)
	form.Append(switchTitleSubtitle(autoDownload, "auto-download", "Auto download videos", "Download new videos."))

	var downloadPolicy data.DownloadPolicy
	if dps, ok := rdx.GetLastVal(data.ChannelDownloadPolicyProperty, channelId); ok && dps != "" {
		downloadPolicy = data.ParseDownloadPolicy(dps)
	}
	form.Append(downloadPolicySelect(downloadPolicy))

	body.Append(submitButton("Update", "manage-channel"))

	if err = strom.WriteResponse(w, root); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func switchTitleSubtitle(on bool, name, title, subTitle string) strom.Element {

	row := strom.Create("ul", atoms.FlexRow(sizes.Normal)...).
		AddAtom(atoms.AlignItemsCenter)

	switchElement := strom.Create("input").
		SetAttribute("id", "name", name).
		SetAttribute("type", "checkbox").
		SetAttribute("switch").
		SetStyle("flex-shrink:0")

	if on {
		switchElement.SetAttribute("checked")
	}

	return row.Append(
		switchElement,
		titleSubtitle(name, title, subTitle))
}

func downloadPolicySelect(currentPolicy data.DownloadPolicy) strom.Element {

	row := strom.Create("ul", atoms.FlexRow(sizes.Normal)...).
		AddAtom(atoms.AlignItemsCenter)

	dps := strom.Create("select").
		SetAttribute("name", "download-policy")

	for _, dp := range data.AllDownloadPolicies() {
		opt := strom.CreateText("option", string(dp))
		if dp == currentPolicy {
			opt.SetAttribute("selected")
		}
		dps.Append(opt)
	}

	row.Append(dps)

	row.Append(titleSubtitle("download-policy", "Download policy", "Recent - limit to the last 10 videos. All - no download limits."))

	return row

	//<select id="download-policy" name="download-policy">
	//	{{$policy := .ChannelDownloadPolicy}}
	//	{{range .AllDownloadPolicies}}
	//	<option {{if eq . $policy}}selected{{end}}>{{.}}</option>
	//	{{end}}
	//	</select>
	//		<label for="download-policy">
	//		<span class="title">Download policy</span>
	//		<span class="subtitle subtle">Recent - limit to the last 10 videos. All - no download limits.</span>
	//		</label>
}

func titleSubtitle(name, title, subTitle string) strom.Element {

	titleStack := strom.Create("label", atoms.FlexCol(sizes.XSmall)...).
		SetAttribute("for", name)

	titleStack.Append(strom.CreateText("span", title).
		AddAtom(atoms.FontWeightBold))
	titleStack.Append(strom.CreateText("span", subTitle).
		SetStyle(
			styles.Decl("font-size", font_sizes.Small),
			styles.Decl("color", colors.Gray)))

	return titleStack
}
