package rest

import (
	"fmt"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/rest/view_models"
	"net/http"
	"time"
)

const (
	recentGroup      = "Less than a week ago"
	thisMonthGroup   = "Less than a month ago"
	thisYearGroup    = "Less than a year ago"
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
	rdx, err = rdx.RefreshWriter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	showAll := r.URL.Query().Has("showAll")

	w.Header().Set("Content-Type", "text/html")

	whKeys := rdx.Keys(data.VideoEndedDateProperty)

	pageTitle := fmt.Sprintf("Last %d watched videos (out of %d)", endedVideosLimit, len(whKeys))
	if showAll {
		pageTitle = fmt.Sprintf("All %d watched videos", len(whKeys))
	}

	hvm := &HistoryViewModel{
		Title:       pageTitle,
		ShowAll:     showAll,
		OpenGroup:   recentGroup,
		GroupsOrder: groupsOrder,
		Groups:      make(map[string][]*view_models.VideoViewModel),
	}

	endedGroups := make(map[string][]string)
	for _, id := range whKeys {
		group := olderGroup
		if ets, ok := rdx.GetLastVal(data.VideoEndedDateProperty, id); ok && ets != "" {
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

		sortedIds, err := rdx.Sort(endedGroups[grp], true, data.VideoEndedDateProperty)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, id := range sortedIds {
			if writtenVideos == endedVideosLimit && !showAll {
				break
			}
			hvm.Groups[grp] = append(hvm.Groups[grp], view_models.GetVideoViewModel(id, rdx, view_models.ShowEndedDate, view_models.ShowPoster))
			writtenVideos++
		}
	}

	if err := tmpl.ExecuteTemplate(w, "history", hvm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
