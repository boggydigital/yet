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
		navButton("Home", "/"),
		navButton("Paste", "/paste"))

	cvs := new(continueVideosSection{rdx: rdx})
	body.Append(strom.OnDemand(cvs.getSectionVideos))

	dvs := new(downloadedVideosSection{rdx: rdx})
	body.Append(strom.OnDemand(dvs.getSectionVideos))

	chs := new(channelsSection{rdx: rdx})
	body.Append(strom.OnDemand(chs.getNewChannels))

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
			if !yield(strom.CreateText("h2", "Continue")) {
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
			if !yield(strom.CreateText("h2", "Videos")) {
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

type channelsSection struct {
	rdx redux.Readable
}

func (cs *channelsSection) getNewChannels() iter.Seq[strom.Element] {
	return func(yield func(element strom.Element) bool) {

		if rdx.Len(data.ChannelAutoRefreshProperty) == 0 {
			return
		}

		var newChannelIds []string

		for channelId := range rdx.Keys(data.ChannelAutoRefreshProperty) {
			if newVideos := yeti.ChannelNotEndedVideos(channelId, math.MaxInt, rdx); len(newVideos) > 0 {
				newChannelIds = append(newChannelIds, channelId)
			}
		}

		if len(newChannelIds) > 0 {
			if !yield(strom.CreateText("h2", "Channels")) {
				return
			}

			channelsContainer := strom.Create("ul", atoms.FlexRowWrap(sizes.Normal)...)

			slices.Sort(newChannelIds)

			cl := new(channelsList{channelIds: newChannelIds, rdx: rdx})
			channelsContainer.Append(strom.OnDemand(cl.getChannelTiles))

			if !yield(channelsContainer) {
				return
			}

		}

	}
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
