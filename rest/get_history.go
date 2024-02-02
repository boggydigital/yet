package rest

import (
	"fmt"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/rest/view_models"
	"net/http"
	"time"
)

const (
	recentGroup      = "A week or less ago"
	thisMonthGroup   = "More than a week, less than a month ago"
	thisYearGroup    = "More than a month, less than a year ago"
	olderGroup       = "More than a year ago"
	endedVideosLimit = 100
)

var groupsOrder = []string{recentGroup, thisMonthGroup, thisYearGroup, olderGroup}

type HistoryViewModel struct {
	Title       string
	ShowAll     bool
	OpenGroup   string
	GroupsOrder []string
	Groups      map[string][]*view_models.VideoViewModel
}

func GetHistory(w http.ResponseWriter, r *http.Request) {

	var err error
	rdx, err = rdx.RefreshReader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	showAll := r.URL.Query().Has("showAll")

	w.Header().Set("Content-Type", "text/html")

	pageTitle := fmt.Sprintf("Last %d watched videos", endedVideosLimit)
	if showAll {
		pageTitle = "All watched videos"
	}

	hvm := &HistoryViewModel{
		Title:       pageTitle,
		ShowAll:     showAll,
		OpenGroup:   recentGroup,
		GroupsOrder: groupsOrder,
		Groups:      make(map[string][]*view_models.VideoViewModel),
	}

	whKeys := rdx.Keys(data.VideoEndedProperty)

	endedGroups := make(map[string][]string)
	for _, id := range whKeys {
		group := olderGroup
		if ets, ok := rdx.GetLastVal(data.VideoEndedProperty, id); ok && ets != "" {
			if et, err := time.Parse(time.RFC3339, ets); err == nil {
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

	writtenVideos := 0

	for _, grp := range groupsOrder {

		if writtenVideos == endedVideosLimit && !showAll {
			break
		}

		if len(endedGroups[grp]) == 0 {
			continue
		}

		sortedIds, err := rdx.Sort(endedGroups[grp], true, data.VideoEndedProperty)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, id := range sortedIds {
			if writtenVideos == endedVideosLimit && !showAll {
				break
			}
			hvm.Groups[grp] = append(hvm.Groups[grp], view_models.GetVideoViewModel(id, rdx, view_models.ShowEndedDate))
			writtenVideos++
		}
	}

	if err := tmpl.ExecuteTemplate(w, "history", hvm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
