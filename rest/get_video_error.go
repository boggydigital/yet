package rest

import (
	"net/http"
	"path"

	"github.com/boggydigital/strom"
	"github.com/boggydigital/strom/styles"
	"github.com/boggydigital/strom/vars/atoms"
	"github.com/boggydigital/strom/vars/colors"
	"github.com/boggydigital/strom/vars/sizes"
	"github.com/boggydigital/yet/data"
)

func GetVideoError(w http.ResponseWriter, r *http.Request) {

	// GET /video_error/{videoId}?err

	videoId := r.PathValue("videoId")
	errStr := r.URL.Query().Get("err")

	root, body := strom.RootBody("Video error", atoms.FlexCol(sizes.Normal)...)

	topRow := strom.Create("ul", atoms.FlexRow(sizes.Small)...).AddAtom(atoms.AlignItemsCenter)
	body.Append(topRow)

	topRow.Append(
		navButton("Home", "/"),
		strom.CreateText("h2", "Video error"))

	if videoTitle, ok := rdx.GetLastVal(data.VideoTitleProperty, videoId); ok && videoTitle != "" {
		body.Append(strom.CreateText("h3", videoTitle))
	}

	body.Append(strom.Create("span").Append(
		strom.CreateText("span", "VideoId: ").SetStyle(styles.Decl("color", colors.Gray)),
		strom.CreateText("span", videoId)))

	body.Append(strom.Create("span").Append(
		strom.CreateText("span", "Error: ").SetStyle(styles.Decl("color", colors.Gray)),
		strom.CreateText("span", errStr)))

	originRow := strom.Create("ul", atoms.FlexRowWrap(sizes.Small)...).
		AddAtom(atoms.AlignItemsCenter)
	body.Append(originRow)

	originRow.Append(
		navButton("Watch at origin", "https://www.youtube.com/watch?v="+videoId),
		navButton("Manage video", path.Join("/manage_video", videoId)))

	if err := strom.WriteResponse(w, root); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
