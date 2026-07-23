package rest

import (
	"iter"
	"maps"
	"math"
	"net/http"
	"slices"

	"github.com/boggydigital/nod"
	"github.com/boggydigital/redux"
	"github.com/boggydigital/strom"
	"github.com/boggydigital/strom/vars/atoms"
	"github.com/boggydigital/strom/vars/sizes"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
)

var jumpToSections = []string{"videos", "channels", "playlists", "history"}
var jumpToSectionTitles = map[string]string{
	"videos":    "Videos",
	"channels":  "Channels",
	"playlists": "Playlists",
	"history":   "History",
}

func GetList(w http.ResponseWriter, r *http.Request) {

	// GET /list

	var err error
	rdx, err = rdx.RefreshWriter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	root, body := strom.RootBody("Watch list", atoms.FlexCol(sizes.Normal)...)

	topRow := strom.Create("ul", atoms.FlexRowWrap(sizes.Small)...).AddAtom(atoms.AlignItemsCenter)
	body.Append(topRow)

	topRow.Append(
		navButton("Search", "/search"),
		navButton("Paste", "/paste"))

	body.Append(strom.CreateText("h2", "Jump to"))

	jumpContainer := strom.Create("ul", atoms.FlexRow(sizes.Small)...)
	body.Append(jumpContainer)

	for _, section := range jumpToSections {
		jumpContainer.Append(navButton(jumpToSectionTitles[section], "#"+section))
	}

	cvs := new(continueVideosSection{rdx: rdx})
	body.Append(strom.OnDemand(cvs.getSectionVideos))

	dvs := new(downloadedVideosSection{rdx: rdx})
	body.Append(strom.OnDemand(dvs.getSectionVideos))

	chs := new(channelsSection{rdx: rdx})
	body.Append(strom.OnDemand(chs.getNewChannels))

	pls := new(playlistsSection{rdx: rdx})
	body.Append(strom.OnDemand(pls.getNewPlaylists))

	body.Append(strom.OnDemand(chs.getCompletedChannels))

	body.Append(strom.OnDemand(pls.getCompletedPlaylists))

	qvs := new(queueudVideosSection{rdx: rdx})
	body.Append(strom.OnDemand(qvs.getSectionVideos))

	body.Append(strom.CreateText("h2", "History").
		SetAttribute("id", "history"))

	body.Append(navButton("See full watch history", "/history"))

	if err = strom.WriteResponse(w, root); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type videosList struct {
	videoIds []string
	rdx      redux.Readable
}

func (vl *videosList) getVideoTiles() iter.Seq[strom.Element] {
	return func(yield func(element strom.Element) bool) {
		for _, videoId := range vl.videoIds {
			if !yield(videoTile(videoId, vl.rdx)) {
				return
			}
		}
	}
}

type continueVideosSection struct {
	rdx redux.Readable
}

func (cvs *continueVideosSection) getSectionVideos() iter.Seq[strom.Element] {
	return func(yield func(element strom.Element) bool) {

		continueVideoIds, err := getContinueVideos(rdx)
		if err != nil {
			nod.Log(err.Error())
			return
		}

		if len(continueVideoIds) > 0 {
			if !yield(strom.CreateText("h2", "Continue").
				SetAttribute("id", "continue")) {
				return
			}

			videosContainer := strom.Create("ul", atoms.FlexRowWrap(sizes.Normal)...)

			vl := new(videosList{videoIds: continueVideoIds, rdx: rdx})
			videosContainer.Append(strom.OnDemand(vl.getVideoTiles))

			if !yield(videosContainer) {
				return
			}
		}
	}
}

type downloadedVideosSection struct {
	rdx redux.Readable
}

func (dvs *downloadedVideosSection) getSectionVideos() iter.Seq[strom.Element] {
	return func(yield func(element strom.Element) bool) {

		downloadedVideoIds, err := getDownloadedVideos(rdx)
		if err != nil {
			nod.Log(err.Error())
			return
		}

		if len(downloadedVideoIds) > 0 {
			if !yield(strom.CreateText("h2", "Videos").
				SetAttribute("id", "videos")) {
				return
			}

			videosContainer := strom.Create("ul", atoms.FlexRowWrap(sizes.Normal)...)

			vl := new(videosList{videoIds: downloadedVideoIds, rdx: rdx})
			videosContainer.Append(strom.OnDemand(vl.getVideoTiles))

			if !yield(videosContainer) {
				return
			}
		}
	}
}

type queueudVideosSection struct {
	rdx redux.Readable
}

func (qvs *queueudVideosSection) getSectionVideos() iter.Seq[strom.Element] {
	return func(yield func(element strom.Element) bool) {

		queuedVideoIds, err := getQueuedDownloads(rdx)
		if err != nil {
			nod.Log(err.Error())
			return
		}

		if len(queuedVideoIds) > 0 {
			if !yield(strom.CreateText("h2", "Queued downloads").
				SetAttribute("id", "downloads")) {
				return
			}

			videosContainer := strom.Create("ul", atoms.FlexRowWrap(sizes.Normal)...)

			vl := new(videosList{videoIds: queuedVideoIds, rdx: rdx})
			videosContainer.Append(strom.OnDemand(vl.getVideoTiles))

			if !yield(videosContainer) {
				return
			}
		}
	}
}

type channelsList struct {
	channelIds []string
	rdx        redux.Readable
}

func (cl *channelsList) getChannelTiles() iter.Seq[strom.Element] {
	return func(yield func(element strom.Element) bool) {
		for _, channelId := range cl.channelIds {
			if !yield(channelTile(channelId, cl.rdx)) {
				return
			}
		}
	}
}

type playlistsList struct {
	playlistIds []string
	rdx         redux.Readable
}

func (pl *playlistsList) getPlaylistTiles() iter.Seq[strom.Element] {
	return func(yield func(element strom.Element) bool) {
		for _, playlistId := range pl.playlistIds {
			if !yield(playlistTile(playlistId, pl.rdx)) {
				return
			}
		}
	}
}

type channelsSection struct {
	rdx redux.Readable
}

func (cs *channelsSection) getChannels(ended bool) iter.Seq[strom.Element] {
	return func(yield func(element strom.Element) bool) {

		if rdx.Len(data.ChannelAutoRefreshProperty) == 0 {
			return
		}

		var channelIds []string

		for channelId := range rdx.Keys(data.ChannelAutoRefreshProperty) {

			newVideos := yeti.ChannelNotEndedVideos(channelId, math.MaxInt, rdx)
			switch len(newVideos) {
			case 0:
				if ended {
					channelIds = append(channelIds, channelId)
				}
			default:
				if !ended {
					channelIds = append(channelIds, channelId)
				}
			}
		}

		sectionTitle := "Channels"
		sectionId := "channels"
		if ended {
			sectionTitle = "Completed channels"
			sectionId = "completed_channels"
		}

		if len(channelIds) > 0 {
			if !yield(strom.CreateText("h2", sectionTitle).
				SetAttribute("id", sectionId)) {
				return
			}

			channelsContainer := strom.Create("ul", atoms.FlexRowWrap(sizes.Normal)...)

			slices.Sort(channelIds)

			cl := new(channelsList{channelIds: channelIds, rdx: rdx})
			channelsContainer.Append(strom.OnDemand(cl.getChannelTiles))

			if !yield(channelsContainer) {
				return
			}
		}
	}
}

func (cs *channelsSection) getNewChannels() iter.Seq[strom.Element] {
	return cs.getChannels(false)
}

func (cs *channelsSection) getCompletedChannels() iter.Seq[strom.Element] {
	return cs.getChannels(true)
}

type playlistsSection struct {
	rdx redux.Readable
}

func (ps *playlistsSection) getPlaylists(ended bool) iter.Seq[strom.Element] {
	return func(yield func(element strom.Element) bool) {

		if rdx.Len(data.PlaylistAutoRefreshProperty) == 0 {
			return
		}

		var playlistIds []string

		for playlistId := range rdx.Keys(data.PlaylistAutoRefreshProperty) {
			newVideos := yeti.PlaylistNotEndedVideos(playlistId, math.MaxInt, rdx)

			switch len(newVideos) {
			case 0:
				if ended {
					playlistIds = append(playlistIds, playlistId)
				}
			default:
				if !ended {
					playlistIds = append(playlistIds, playlistId)
				}
			}
		}

		sectionTitle := "Playlists"
		sectionId := "playlists"
		if ended {
			sectionTitle = "Completed playlists"
			sectionId = "completed_playlists"
		}

		if len(playlistIds) > 0 {
			if !yield(strom.CreateText("h2", sectionTitle).
				SetAttribute("id", sectionId)) {
				return
			}

			playlistsContainer := strom.Create("ul", atoms.FlexRowWrap(sizes.Normal)...)

			slices.Sort(playlistIds)

			cl := new(playlistsList{playlistIds: playlistIds, rdx: rdx})
			playlistsContainer.Append(strom.OnDemand(cl.getPlaylistTiles))

			if !yield(playlistsContainer) {
				return
			}
		}
	}
}

func (ps *playlistsSection) getNewPlaylists() iter.Seq[strom.Element] {
	return ps.getPlaylists(false)
}

func (ps *playlistsSection) getCompletedPlaylists() iter.Seq[strom.Element] {
	return ps.getPlaylists(true)
}

func getContinueVideos(rdx redux.Readable) ([]string, error) {

	cvs := make(map[string]any)
	var err error

	if rdx.Len(data.VideoProgressProperty) == 0 {
		return nil, nil
	}

	for id := range rdx.Keys(data.VideoProgressProperty) {
		if et, ok := rdx.GetLastVal(data.VideoEndedDateProperty, id); ok && et != "" {
			continue
		}
		cvs[id] = nil
	}

	videoIds := slices.Collect(maps.Keys(cvs))

	if videoIds, err = rdx.Sort(videoIds, false, data.VideoTitleProperty); err == nil {
		return videoIds, nil
	} else {
		return nil, err
	}
}

func getDownloadedVideos(rdx redux.Readable) ([]string, error) {

	dvs := make([]string, 0, rdx.Len(data.VideoDownloadCompletedProperty))

	if rdx.Len(data.VideoDownloadCompletedProperty) == 0 {
		return dvs, nil
	}

	// videos is all downloaded videos that are not:
	// - in history (ended)
	// - in continue (have progress)
	// - is favorite
	// - in any auto-refreshing channel
	// - in any auto-refreshing playlist

	for id := range rdx.Keys(data.VideoDownloadCompletedProperty) {

		if rdx.HasKey(data.VideoEndedDateProperty, id) {
			continue
		}
		if rdx.HasKey(data.VideoProgressProperty, id) {
			continue
		}
		if rdx.HasKey(data.VideoFavoriteProperty, id) {
			continue
		}

		// check if this video is an auto-refreshing channel video
		skip := false
		for channelId := range rdx.Keys(data.ChannelAutoRefreshProperty) {
			if rdx.HasValue(data.ChannelVideosProperty, channelId, id) {
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		// check if this video is an auto-refreshing playlist video
		skip = false
		for playlistId := range rdx.Keys(data.PlaylistAutoRefreshProperty) {
			if rdx.HasValue(data.PlaylistVideosProperty, playlistId, id) {
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		dvs = append(dvs, id)
	}

	var err error
	if dvs, err = rdx.Sort(dvs, false, data.VideoTitleProperty); err == nil {
		return dvs, nil
	} else {
		return nil, err
	}
}

func getQueuedDownloads(rdx redux.Readable) ([]string, error) {

	qdLen := rdx.Len(data.VideoDownloadQueuedProperty)

	qds := make([]string, 0, qdLen)

	if qdLen == 0 {
		return qds, nil
	}

	for id := range rdx.Keys(data.VideoDownloadQueuedProperty) {

		dqTime := ""
		if dqt, ok := rdx.GetLastVal(data.VideoDownloadQueuedProperty, id); ok {
			dqTime = dqt
		}

		// only continue if download was completed _after_ it was queued,
		// meaning it wasn't re-queued again after completion
		if dcd, ok := rdx.GetLastVal(data.VideoDownloadCompletedProperty, id); ok && dcd > dqTime {
			continue
		}

		qds = append(qds, id)
	}

	var err error
	if qds, err = rdx.Sort(qds, false, data.VideoTitleProperty); err == nil {
		return qds, nil
	} else {
		return nil, err
	}
}
