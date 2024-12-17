package rest

import (
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/rest/compton_elements"
	"net/http"
)

func GetCompton(w http.ResponseWriter, r *http.Request) {

	p := compton.Page("compton test area")

	videoIds := []string{
		"yn7kUDRVcVM",
		"ANyJVMhOpkk",
		"jwVEhEPK9dI",
		"KuCRvr6R8Lc",
		"xF8huW3imyk",
	}

	pageStack := compton.FlexItems(p, direction.Column)
	p.Append(pageStack)

	gridItems := compton.GridItems(p)
	pageStack.Append(gridItems)

	for _, videoId := range videoIds {

		videoLink := compton_elements.VideoPoster(p, videoId, rdx)
		gridItems.Append(videoLink)

	}

	if err := p.WriteResponse(w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

}
