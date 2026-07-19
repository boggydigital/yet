package rest

import (
	"math"
	"path"
	"strconv"
	"time"

	"github.com/boggydigital/redux"
	"github.com/boggydigital/strom"
	"github.com/boggydigital/strom/styles"
	"github.com/boggydigital/strom/vars/atoms"
	"github.com/boggydigital/strom/vars/calc"
	"github.com/boggydigital/strom/vars/colors"
	"github.com/boggydigital/strom/vars/font_sizes"
	"github.com/boggydigital/strom/vars/font_weights"
	"github.com/boggydigital/strom/vars/sizes"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
)

var reasonTitles = map[data.VideoEndedReason]string{
	data.Completed:  "Completed",
	data.SeenEnough: "Seen enough",
	data.Skipped:    "Skipped",
}

func videoTile(videoId string, rdx redux.Readable) strom.Element {

	tileContainer := strom.Create("a", atoms.FlexColWrap(sizes.Normal)...).
		SetAttribute("href", path.Join("/watch", videoId)).
		SetStyle(
			"position:relative",
			styles.Decl("width", calc.Mult(sizes.XXXLarge, 1.5)))

	var ended bool
	if rdx.HasKey(data.VideoEndedDateProperty, videoId) {
		ended = true
	}

	poster := strom.Create("img").
		SetAttribute("src", path.Join("/poster?v="+videoId+"&q=hqdefault")).
		SetAttribute("loading", "lazy").
		SetStyle(
			"aspect-ratio:16/9",
			"width:100%",
			"object-fit:cover",
			styles.Decl("border-radius", sizes.XSmall))

	tileContainer.Append(poster)

	if ended {
		poster.SetStyle("filter:grayscale(1.0)")

		reason := data.Completed
		if ver, ok := rdx.GetLastVal(data.VideoEndedReasonProperty, videoId); ok && ver != "" {
			reason = data.ParseVideoEndedReason(ver)
		}

		tileContainer.Append(strom.CreateText("span", reasonTitles[reason]).
			SetStyle(
				"position:absolute",
				"top:0",
				"right:0",
				styles.Decl("font-size", font_sizes.XSmall),
				styles.Decl("padding", sizes.Small),
				styles.Decl("border-bottom-left-radius", sizes.Small),
				styles.Decl("background-color", colors.Background)))
	}

	if durs, sure := rdx.GetLastVal(data.VideoDurationProperty, videoId); sure && durs != "" && durs != "0" {
		if duri, err := strconv.ParseInt(durs, 10, 64); err == nil {

			var remaining int64

			if cts, ok := rdx.GetLastVal(data.VideoProgressProperty, videoId); ok && cts != "" {
				var cti int64
				if cti, err = strconv.ParseInt(cts, 10, 64); err == nil {
					remaining = duri - cti
				}
			}

			durationItems := strom.Create("span", atoms.FlexRowWrap(sizes.Small)...).
				SetStyle(
					"position:absolute",
					"top:0",
					"left:0",
					styles.Decl("font-size", font_sizes.XSmall),
					styles.Decl("padding", sizes.Small),
					styles.Decl("border-bottom-right-radius", sizes.Small),
					styles.Decl("background-color", colors.Background))

			durSpan := strom.CreateText("span", formatSeconds(duri)).
				SetStyle(styles.Decl("font-size", font_sizes.XSmall))

			if remaining > 0 {
				remSpan := strom.CreateText("span", formatSeconds(remaining), atoms.FontWeightBold)
				durationItems.Append(remSpan)
				durSpan.SetStyle(styles.Decl("color", colors.Gray))
			} else {
				if !ended {
					durSpan.AddAtom(atoms.FontWeightBold)
				}
			}

			durationItems.Append(durSpan)

			tileContainer.Append(durationItems)
		}
	}

	titlePropertiesStack := strom.Create("ul", atoms.FlexColWrap(sizes.Small)...)

	if title, ok := rdx.GetLastVal(data.VideoTitleProperty, videoId); ok && title != "" {
		titlePropertiesStack.Append(strom.CreateText("h3", title))
	}

	vsp := videoSummaryProperties(videoId, rdx)

	propertiesStack := strom.Create("ul", atoms.DisplayFlex, atoms.FlexDirColumn).
		SetStyle(styles.Decl("row-gap", sizes.XSmall))
	titlePropertiesStack.Append(propertiesStack)

	for _, p := range propertiesOrder {
		v := vsp[p]
		if v == "" {
			continue
		}

		ptv := propertyTitles[p] + ": " + v

		propertyRow := strom.CreateText("span", ptv).
			SetStyle(
				styles.Decl("color", colors.Gray),
				styles.Decl("font-size", font_sizes.XSmall))

		propertiesStack.Append(propertyRow)
	}

	tileContainer.Append(titlePropertiesStack)

	return tileContainer
}

func videoSummaryProperties(videoId string, rdx redux.Readable) map[string]string {
	properties := make(map[string]string)

	if och, ok := rdx.GetLastVal(data.VideoOwnerChannelNameProperty, videoId); ok && och != "" {
		properties[data.VideoOwnerChannelNameProperty] = och
	}

	if ets, ok := rdx.GetLastVal(data.VideoEndedDateProperty, videoId); ok && ets != "" {
		properties[data.VideoEndedDateProperty] = parseAndFormatDate(ets)
	}

	if len(properties) < 2 {
		var publishedDate string
		if pds, ok := rdx.GetLastVal(data.VideoPublishDateProperty, videoId); ok && pds != "" {
			publishedDate = parseAndFormatDate(pds)
		} else {
			if ptts, sure := rdx.GetLastVal(data.VideoPublishTimeTextProperty, videoId); sure && ptts != "" {
				publishedDate = ptts
			}
		}

		if publishedDate != "" {
			properties[data.VideoPublishDateProperty] = publishedDate
		}
	}

	if len(properties) < 2 {
		if dts, ok := rdx.GetLastVal(data.VideoDownloadCompletedProperty, videoId); ok && dts != "" {
			properties[data.VideoDownloadCompletedProperty] = parseAndFormatDate(dts)
		}
	}

	return properties
}

func parseAndFormat(ts string) string {
	if pt, err := time.Parse(time.RFC3339, ts); err == nil {
		return pt.Local().Format(time.RFC1123)
	} else {
		return ts
	}
}

func parseAndFormatDate(ts string) string {
	if pt, err := time.Parse(time.RFC3339, ts); err == nil {
		return pt.Local().Format("Mon, 2 Jan 2006")
	} else {
		return ts
	}
}

func formatSeconds(ts int64) string {
	if ts == 0 {
		return "unknown"
	}

	t := time.Unix(ts, 0).UTC()

	layout := "4:05"
	if t.Hour() > 0 {
		layout = "15:04:05"
	}

	return t.Format(layout)
}

func channelTile(channelId string, rdx redux.Readable) strom.Element {

	var channelTitle string
	if tp, ok := rdx.GetLastVal(data.ChannelTitleProperty, channelId); ok && tp != "" {
		channelTitle = tp
	}

	newVideos := len(yeti.ChannelNotEndedVideos(channelId, math.MaxInt, rdx))

	return linkTile(path.Join("/channel", channelId), newVideos, channelTitle)
}

func playlistTile(playlistId string, rdx redux.Readable) strom.Element {

	var playlistTitle string
	if tp, ok := rdx.GetLastVal(data.PlaylistTitleProperty, playlistId); ok && tp != "" {
		playlistTitle = tp
	}

	var channelTitle string
	if channel, ok := rdx.GetLastVal(data.PlaylistChannelProperty, playlistId); ok && channel != "" {
		channelTitle = channel
	}

	newVideos := len(yeti.PlaylistNotEndedVideos(playlistId, math.MaxInt, rdx))

	return linkTile(path.Join("/playlist", playlistId), newVideos, playlistTitle, channelTitle)
}

func linkTile(href string, count int, titles ...string) strom.Element {

	tileContainer := strom.Create("a", atoms.FlexRow(sizes.Small)...).
		SetAttribute("href", href).
		AddAtom(atoms.AlignItemsCenter, atoms.BorderRadiusSmall, atoms.PaddingSmall).
		SetStyle(
			"flow-shrink:0",
			"width:fit-content",
			styles.Decl("padding-inline-end", sizes.Normal),
			styles.Decl("background", colors.Highlight))

	if count > 0 {
		tileContainer.Append(strom.CreateText("span", strconv.Itoa(count)).
			SetStyle(
				styles.Decl("border-radius", sizes.Small),
				styles.Decl("padding", sizes.Small),
				styles.Decl("background-color", colors.Background),
				styles.Decl("font-weight", font_weights.Bold),
				styles.Decl("font-size", font_sizes.XXSmall)))
	}

	titlesStack := strom.Create("ul", atoms.FlexCol(sizes.XSmall)...)
	tileContainer.Append(titlesStack)

	if len(titles) > 0 {
		titlesStack.Append(strom.CreateText("span", titles[0], atoms.FontWeightBold))
	}

	if len(titles) > 1 {
		titlesStack.Append(strom.CreateText("span", titles[1]).
			SetStyle(
				styles.Decl("font-size", font_sizes.XSmall),
				styles.Decl("color", colors.Gray)))
	}

	return tileContainer
}

func navButton(title, href string) strom.Element {
	return strom.Create("a").
		SetTextContent(title).
		SetAttribute("href", href).
		SetStyle(buttonStyles()...)
}

func submitButton(value, form string) strom.Element {
	return strom.Create("input").
		SetAttribute("type", "submit").
		SetAttribute("form", form).
		SetAttribute("value", value).
		SetStyle("appearance:none").
		SetStyle(buttonStyles()...)
}

func buttonStyles() []string {
	return []string{
		"border:none",
		"width:fit-content",
		styles.Decl("padding-block", sizes.Small),
		styles.Decl("padding-inline", calc.Mult(sizes.Small, 1.25)),
		styles.Decl("background-color", colors.Highlight),
		styles.Decl("border-radius", sizes.Small),
		styles.Decl("color", colors.Foreground),
		styles.Decl("font-size", font_sizes.Normal),
	}
}

func textInputStyles() []string {
	return []string{
		"appearance:none",
		"border:none",
		styles.Decl("border-radius", sizes.Small),
		styles.Decl("max-width", calc.Mult(sizes.XXXLarge, 1.5)),
		styles.Decl("padding", sizes.Small),
		styles.Decl("font-size", font_sizes.Normal),
	}
}
