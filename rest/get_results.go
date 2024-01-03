package rest

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yet/rest/view_models"
	"github.com/boggydigital/yt_urls"
	"net/http"
	"strings"
)

type ResultsViewModel struct {
	SearchQuery string
	Refinements []string
	Videos      []*view_models.VideoViewModel
}

func GetResults(w http.ResponseWriter, r *http.Request) {

	searchQuery := r.URL.Query().Get("search_query")

	terms := strings.Split(searchQuery, " ")

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

	rvm := &ResultsViewModel{
		SearchQuery: searchQuery,
		Refinements: sid.Refinements,
	}

	for _, vr := range sid.VideoRenderers() {
		rvm.Videos = append(rvm.Videos, view_models.GetVideoViewModel(vr.VideoId, vmRdx,
			view_models.ShowOwnerChannel,
			view_models.ShowPublishedDate,
			view_models.ShowViewCount))
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
