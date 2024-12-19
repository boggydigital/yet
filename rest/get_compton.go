package rest

import (
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/rest/compton_elements"
	"net/http"
	"slices"
)

func GetCompton(w http.ResponseWriter, r *http.Request) {

	p := compton.Page("compton test area")

	videoIds := rdx.Keys(data.VideoTitleProperty)
	slices.Sort(videoIds)
	videoIds = videoIds[:20]

	pageStack := compton.FlexItems(p, direction.Column)
	p.Append(pageStack)

	gridItems := compton.GridItems(p)
	pageStack.Append(gridItems)

	for _, videoId := range videoIds {
		videoLink := compton_elements.VideoLink(p, videoId, rdx)
		gridItems.Append(videoLink)

	}

	if err := p.WriteResponse(w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

}
