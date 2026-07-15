package rest

import (
	"iter"
	"net/http"

	"github.com/boggydigital/redux"
	"github.com/boggydigital/strom"
	"github.com/boggydigital/strom/vars"
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
		url := "/refresh_channel_videos?id=" + channelId
		http.Redirect(w, r, url, http.StatusPermanentRedirect)
		return
	}

	var channelTitle string
	if ct, ok := rdx.GetLastVal(data.ChannelTitleProperty, channelId); ok && ct != "" {
		channelTitle = ct
	}

	root := strom.Page(channelTitle)

	var body strom.Element
	for body = range root.GetElementsByTagName("body") {
		break
	}

	body.AddClass("d-f", "fd-c", "rg-l")

	body.Append(navButton("Home", "/"))

	body.Append(channelTile(channelId, rdx))

	if cd, ok := rdx.GetLastVal(data.ChannelDescriptionProperty, channelId); ok && cd != "" {
		body.Append(strom.CreateText("span", cd).
			SetStyle(map[string]string{
				"color":     vars.Color(vars.ColorGray),
				"max-width": "calc(4 * " + vars.Size(vars.SizeXXXLarge) + ")",
			}))
	}

	channelNavButtonsRow := strom.Create("ul", "d-f", "fd-r", "cg-n", "rg-n").
		Append(navButton("RSS", "https://www.youtube.com/feeds/videos.xml?channel_id="+channelId)).
		Append(navButton("Playlists", "/channel_playlists?id="+channelId)).
		Append(navButton("Refresh", "/refresh_channel_videos?id="+channelId)).
		Append(navButton("Manage", "/manage_channel?id="+channelId))

	body.Append(channelNavButtonsRow)

	cv := new(channelVideos{channelId: channelId, rdx: rdx})

	videos := strom.Create("ul", "d-f", "cg-l", "rg-l").
		SetStyle(map[string]string{
			"flex-flow": "row wrap",
		})
	body.Append(videos)

	videos.Append(strom.Defer(cv.getVideos))

	if err = strom.WriteResponse(w, root); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	//w.Header().Set("Content-Type", "text/html")
	//
	//if err := tmpl.ExecuteTemplate(w, "channels", view_models.GetChannelViewModel(channelId, rdx)); err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
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
