package compton_elements

import (
	"embed"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/size"
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/yet/data"
)

//go:embed "styles/*.css"
var videoLinkStyles embed.FS

func VideoPoster(r compton.Registrar, videoId string, rdx kevlar.ReadableRedux) compton.Element {

	r.RegisterStyles(videoLinkStyles, "styles/video-link.css")

	var dehydratedPoster string
	var repColor string

	if dhp, ok := rdx.GetLastVal(data.VideoDehydratedPosterProperty, videoId); ok {
		dehydratedPoster = dhp
	}

	if rc, ok := rdx.GetLastVal(data.VideoDehydratedRepColorProperty, videoId); ok {
		repColor = rc
	}

	issaImage := compton.IssaImageDehydrated(r, repColor, dehydratedPoster, "/poster?v="+videoId+"&q=hqdefault")
	issaImage.AddClass("video-poster")

	issaImage.Width(size.ColumnWidth).AspectRatio(float64(16) / float64(9))

	return issaImage
}
