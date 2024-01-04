package rest

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/pathology"
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
	Channels    []*view_models.ChannelViewModel
	Playlists   []*view_models.PlaylistViewModel
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

	metadataDir, err := pathology.GetAbsDir(paths.Metadata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	extractedProperties := extractedSearchVideosProperties
	extractedProperties = append(extractedProperties, extractedSearchPlaylistProperties...)
	extractedProperties = append(extractedProperties, extractedSearchChannelProperties...)

	wRdx, err := kvas.NewReduxWriter(metadataDir, extractedProperties...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	propertyValues := extractSearchVideosMetadata(sid.VideoRenderers())
	for property, keyValues := range propertyValues {
		if err := wRdx.BatchAddValues(property, keyValues); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	propertyValues = extractSearchPlaylistMetadata(sid.PlaylistRenderers())
	for property, keyValues := range propertyValues {
		if err := wRdx.BatchAddValues(property, keyValues); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	propertyValues = extractSearchChannelMetadata(sid.ChannelRenderers())
	for property, keyValues := range propertyValues {
		if err := wRdx.BatchAddValues(property, keyValues); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	rvm := &ResultsViewModel{
		SearchQuery: searchQuery,
		Refinements: sid.Refinements,
	}

	for _, chr := range sid.ChannelRenderers() {
		rvm.Channels = append(rvm.Channels, view_models.GetChannelViewModel(chr.ChannelId, wRdx))
	}

	for _, plr := range sid.PlaylistRenderers() {
		rvm.Playlists = append(rvm.Playlists, view_models.GetPlaylistViewModel(plr.PlaylistId, wRdx))
	}

	for _, vr := range sid.VideoRenderers() {
		rvm.Videos = append(rvm.Videos, view_models.GetVideoViewModel(vr.VideoId, wRdx,
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

var extractedSearchVideosProperties = []string{
	data.VideoTitleProperty,
	data.VideoOwnerChannelNameProperty,
	data.VideoViewCountProperty,
	data.VideoPublishTimeTextProperty,
	data.VideoEndedProperty,
}

var extractedSearchPlaylistProperties = []string{
	data.PlaylistTitleProperty,
	data.PlaylistChannelProperty,
}

var extractedSearchChannelProperties = []string{
	data.ChannelTitleProperty,
	data.ChannelDescriptionProperty,
}

func extractSearchVideosMetadata(svrs []yt_urls.VideoRenderer) map[string]map[string][]string {
	pkv := make(map[string]map[string][]string)

	for _, property := range extractedSearchVideosProperties {

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

func extractSearchPlaylistMetadata(sprs []yt_urls.PlaylistRenderer) map[string]map[string][]string {
	pkv := make(map[string]map[string][]string)

	for _, property := range extractedSearchPlaylistProperties {

		pkv[property] = make(map[string][]string)

		for _, pvr := range sprs {

			id := pvr.PlaylistId

			switch property {
			case data.PlaylistTitleProperty:
				pkv[property][id] = []string{pvr.Title.SimpleText}
			}
		}

	}

	return pkv
}

func extractSearchChannelMetadata(scrs []yt_urls.ChannelRenderer) map[string]map[string][]string {
	pkv := make(map[string]map[string][]string)

	for _, property := range extractedSearchChannelProperties {

		pkv[property] = make(map[string][]string)

		for _, cr := range scrs {

			id := cr.ChannelId

			switch property {
			case data.ChannelTitleProperty:
				pkv[property][id] = []string{cr.Title.SimpleText}
			case data.ChannelDescriptionProperty:
				pkv[property][id] = []string{cr.DescriptionSnippet.String()}
			}
		}

	}

	return pkv
}
