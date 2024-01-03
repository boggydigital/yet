package rest

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yt_urls"
	"net/http"
	"strings"
)

type ResultsViewModel struct {
	Videos []*VideoViewModel
}

func GetResults(w http.ResponseWriter, r *http.Request) {

	terms := strings.Split(r.URL.Query().Get("search_query"), " ")

	sid, err := yt_urls.GetSearchResultsPage(http.DefaultClient, terms...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	metadataDir, err := paths.GetAbsDir(paths.Metadata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	vmRdx, err := kvas.NewReduxWriter(metadataDir, extractedSearchProperties...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	properyValues := extractSearchVideosMetadata(sid.VideoRenderers())
	for property, keyValues := range properyValues {
		if err := vmRdx.BatchAddValues(property, keyValues); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	rvm := &ResultsViewModel{}

	for _, vr := range sid.VideoRenderers() {
		rvm.Videos = append(rvm.Videos, videoViewModel(vr.VideoId, vmRdx,
			ShowOwnerChannel,
			ShowPublishedDate,
			ShowViewCount))
	}

	w.Header().Set("Content-Type", "text/html")

	if err := tmpl.ExecuteTemplate(w, "results", rvm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

var extractedSearchProperties = []string{
	data.VideoTitleProperty,
	data.VideoOwnerChannelNameProperty,
	data.VideoViewCountProperty,
	data.VideoPublishTimeTextProperty,
	data.VideoEndedProperty,
}

func extractSearchVideosMetadata(svrs []yt_urls.VideoRenderer) map[string]map[string][]string {
	pkv := make(map[string]map[string][]string)

	for _, property := range extractedSearchProperties {

		pkv[property] = make(map[string][]string)

		for _, svr := range svrs {

			id := svr.VideoId

			switch property {
			case data.VideoTitleProperty:
				pkv[property][id] = []string{svr.Title.String()}
			case data.VideoOwnerChannelNameProperty:
				pkv[property][id] = []string{svr.OwnerText.String()}
			case data.VideoViewCountProperty:
				pkv[property][id] = []string{svr.ViewCountText.SimpleText}
			case data.VideoPublishTimeTextProperty:
				pkv[property][id] = []string{svr.PublishedTimeText.SimpleText}
			}
		}

	}

	return pkv
}
