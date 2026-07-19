package rest

import (
	"fmt"
	"iter"
	"maps"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/boggydigital/nod"
	"github.com/boggydigital/redux"
	"github.com/boggydigital/strom"
	"github.com/boggydigital/strom/styles"
	"github.com/boggydigital/strom/vars/atoms"
	"github.com/boggydigital/strom/vars/colors"
	"github.com/boggydigital/strom/vars/sizes"
	"github.com/boggydigital/yet/data"
)

const (
	recentGroup      = "Less than a week ago"
	thisMonthGroup   = "Less than a month ago"
	thisYearGroup    = "Less than a year ago"
	olderGroup       = "More than a year ago"
	endedVideosLimit = 100
)

var groupsOrder = []string{recentGroup, thisMonthGroup, thisYearGroup, olderGroup}

func GetHistory(w http.ResponseWriter, r *http.Request) {

	var err error
	rdx, err = rdx.RefreshWriter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	showAll := r.URL.Query().Has("showAll")

	whKeys := rdx.Keys(data.VideoEndedDateProperty)
	whKeysLen := rdx.Len(data.VideoEndedDateProperty)

	endedGroups := make(map[string][]string)
	for id := range whKeys {
		group := olderGroup
		if ets, ok := rdx.GetLastVal(data.VideoEndedDateProperty, id); ok && ets != "" {
			var et time.Time
			if et, err = time.Parse(time.RFC3339, ets); err == nil {
				days := time.Now().Sub(et).Hours() / 24
				if days <= 7 {
					group = recentGroup
				} else if days <= 30 {
					group = thisMonthGroup
				} else if days <= 365 {
					group = thisYearGroup
				}
			}
		}
		endedGroups[group] = append(endedGroups[group], id)
	}

	pageTitle := fmt.Sprintf("Last %d watched videos", endedVideosLimit)
	if showAll {
		pageTitle = fmt.Sprintf("All %d watched videos", whKeysLen)
	}

	root, body := strom.RootBody(pageTitle, atoms.FlexCol(sizes.Normal)...)

	topRow := strom.Create("ul", atoms.FlexRow(sizes.Small)...).AddAtom(atoms.AlignItemsCenter)
	body.Append(topRow)

	topRow.Append(navButton("Home", "/"))
	topRow.Append(strom.CreateText("h2", pageTitle))

	hs := new(historyStats{rdx: rdx})
	body.Append(strom.Create("ul", atoms.FlexRowWrap(sizes.Normal)...).
		SetStyle(styles.Decl("color", colors.Gray)).
		Append(strom.OnDemand(hs.getStats)))

	egv := new(endedGroupsVideos{endedGroups: endedGroups, showAll: showAll, rdx: rdx})
	body.Append(strom.OnDemand(egv.getVideoSections))

	if !showAll {
		body.Append(strom.Create("br"))
		body.Append(strom.CreateText("span", "To load this page faster, yet is limiting displayed videos.").
			SetStyle(styles.Decl("color", colors.Gray)))
		body.Append(navButton("Show all", "/history?showAll"))
	}

	if err = strom.WriteResponse(w, root); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type endedGroupsVideos struct {
	endedGroups map[string][]string
	showAll     bool
	rdx         redux.Readable
}

func (egv *endedGroupsVideos) getVideoSections() iter.Seq[strom.Element] {
	return func(yield func(strom.Element) bool) {

		writtenVideos := 0

		for _, grp := range groupsOrder {

			if writtenVideos == endedVideosLimit && !egv.showAll {
				return
			}

			if len(egv.endedGroups[grp]) == 0 {
				continue
			}

			sortedIds, err := rdx.Sort(egv.endedGroups[grp], true, data.VideoEndedDateProperty)
			if err != nil {
				nod.Log(err.Error())
				return
			}

			if !yield(strom.CreateText("h2", grp)) {
				return
			}

			videosContainer := strom.Create("ul", atoms.FlexRowWrap(sizes.Normal)...)

			switch egv.showAll {
			case true:
				vl := new(videosList{sortedIds, rdx})
				videosContainer.Append(strom.OnDemand(vl.getVideoTiles))
			case false:
				lv := new(limitedVideos{sortedIds, &writtenVideos, rdx})
				videosContainer.Append(strom.OnDemand(lv.getVideoTiles))
			}

			if !yield(videosContainer) {
				return
			}
		}

	}
}

type limitedVideos struct {
	videoIds      []string
	writtenVideos *int
	rdx           redux.Readable
}

func (lv *limitedVideos) getVideoTiles() iter.Seq[strom.Element] {
	return func(yield func(element strom.Element) bool) {
		for _, videoId := range lv.videoIds {
			if *lv.writtenVideos == endedVideosLimit {
				return
			}
			if !yield(videoTile(videoId, lv.rdx)) {
				return
			}
			*lv.writtenVideos++
		}
	}
}

type historyStats struct {
	rdx redux.Readable
}

func (hs *historyStats) getStats() iter.Seq[strom.Element] {
	return func(yield func(element strom.Element) bool) {

		stats := make(map[data.VideoEndedReason]int)

		for videoId := range rdx.Keys(data.VideoEndedDateProperty) {
			// do not check presense as an empty value indicated default (completed)
			ver, _ := rdx.GetLastVal(data.VideoEndedReasonProperty, videoId)
			switch ver {
			case "":
				stats[data.Completed]++
			default:
				videoEndedReason := data.ParseVideoEndedReason(ver)
				stats[videoEndedReason]++
			}

		}

		totalEnded := rdx.Len(data.VideoEndedDateProperty)

		if !yield(strom.CreateText("span", "Total ended: "+strconv.Itoa(totalEnded))) {
			return
		}

		erKeys := slices.Collect(maps.Keys(stats))
		slices.Sort(erKeys)

		for _, er := range erKeys {
			if !yield(strom.CreateText("span", reasonTitles[er]+": "+strconv.Itoa(stats[er]))) {
				return
			}
		}
	}
}
