package rest

import (
	"fmt"
	"github.com/boggydigital/yet/data"
	"io"
	"net/http"
	"strings"
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

func GetHistory(w http.ResponseWriter, r *http.Request) {

	showAll := r.URL.Query().Has("showAll")

	w.Header().Set("Content-Type", "text/html")

	sb := &strings.Builder{}
	sb.WriteString("<!doctype html>")
	sb.WriteString("<html>")
	sb.WriteString("<head>" +
		"<meta charset='UTF-8'>" +
		"<link rel='icon' href='data:image/svg+xml,<svg xmlns=%22http://www.w3.org/2000/svg%22 viewBox=%220 0 100 100%22><text y=%22.9em%22 font-size=%2290%22>🔻</text></svg>' type='image/svg+xml'/>" +
		"<meta name='viewport' content='width=device-width, initial-scale=1.0'>" +
		"<meta name='color-scheme' content='dark light'>" +
		"<title>🔻 History</title>" +
		"<style>")

	writeSharedStyles(sb)

	// no history specific styles at the moment
	sb.WriteString("a.video.showAll {color: dodgerblue}")

	sb.WriteString("</style></head>")
	sb.WriteString("<body>")

	pageTitle := fmt.Sprintf("Last %d watched videos", endedVideosLimit)
	if showAll {
		pageTitle = "All watched videos"
	}

	sb.WriteString("<h1>" + pageTitle + "</h1>")

	whKeys := rdx.Keys(data.VideoEndedProperty)

	endedGroups := make(map[string][]string)
	for _, id := range whKeys {
		group := olderGroup
		if ets, ok := rdx.GetFirstVal(data.VideoEndedProperty, id); ok && ets != "" {
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

		open := ""
		if grp == recentGroup {
			open = "open"
		}

		sb.WriteString("<details " + open + "><summary><h2>" + grp + "</h2></summary>")

		sortedIds, err := rdx.Sort(endedGroups[grp], true, data.VideoEndedProperty)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, id := range sortedIds {
			if writtenVideos == endedVideosLimit && !showAll {
				break
			}
			writeVideo(id, rdx, sb, ShowEndedDate)
			writtenVideos++
		}
		sb.WriteString("</details>")
	}

	if !showAll {
		sb.WriteString("<div class='subtle'>" +
			"To load this page faster, yet is limiting displayed videos. " +
			"Click the link below to see all videos you ever watched on yet" +
			"</div>")
		sb.WriteString("<a class='video showAll' href='/history?showAll'>Show all videos</a>")
	}

	sb.WriteString("</body>")
	sb.WriteString("</html>")

	if _, err := io.WriteString(w, sb.String()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
