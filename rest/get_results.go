package rest

import (
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

	rvm := &ResultsViewModel{}

	for _, content := range sid.Contents.TwoColumnSearchResultsRenderer.PrimaryContents.SectionListRenderer.Contents {
		for _, isrc := range content.ItemSectionRenderer.Contents {
			if isrc.VideoRenderer.VideoId == "" {
				continue
			}
			rvm.Videos = append(rvm.Videos, videoViewModel(isrc.VideoRenderer.VideoId, rdx))
		}
	}

	w.Header().Set("Content-Type", "text/html")

	if err := tmpl.ExecuteTemplate(w, "results", rvm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
