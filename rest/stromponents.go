package rest

import (
	"math"
	"path"
	"strconv"
	"time"

	"github.com/boggydigital/redux"
	"github.com/boggydigital/strom"
	"github.com/boggydigital/strom/vars/atoms"
	"github.com/boggydigital/strom/vars/calc"
	"github.com/boggydigital/strom/vars/colors"
	"github.com/boggydigital/strom/vars/font_sizes"
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
		SetStyle(map[string]string{
			"width":    calc.Mult(sizes.XXXLarge, 1.5),
			"position": "relative",
		})

	var ended bool
	if rdx.HasKey(data.VideoEndedDateProperty, videoId) {
		ended = true
	}

	poster := strom.Create("img").
		SetAttribute("src", path.Join("/poster?v="+videoId+"&q=hqdefault")).
		SetAttribute("loading", "lazy").
		SetStyle(map[string]string{
			"border-radius": sizes.XSmall,
			"aspect-ratio":  "16/9",
			"width":         "100%",
			"object-fit":    "cover",
		})

	tileContainer.Append(poster)

	if ended {
		poster.SetStyle(map[string]string{
			"filter": "grayscale(1.0)",
		})

		reason := data.Completed
		if ver, ok := rdx.GetLastVal(data.VideoEndedReasonProperty, videoId); ok && ver != "" {
			reason = data.ParseVideoEndedReason(ver)
		}

		tileContainer.Append(strom.CreateText("span", reasonTitles[reason]).
			SetStyle(map[string]string{
				"position":                  "absolute",
				"top":                       "0",
				"right":                     "0",
				"font-size":                 font_sizes.XSmall,
				"padding":                   sizes.Small,
				"border-bottom-left-radius": sizes.Small,
				"background-color":          colors.Background,
			}))
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
				SetStyle(map[string]string{
					"position":                   "absolute",
					"top":                        "0",
					"left":                       "0",
					"font-size":                  font_sizes.XSmall,
					"padding":                    sizes.Small,
					"border-bottom-right-radius": sizes.Small,
					"background-color":           colors.Background,
				})

			durSpan := strom.CreateText("span", formatSeconds(duri)).
				SetStyle(map[string]string{
					"font-size": font_sizes.XSmall,
				})

			if remaining > 0 {
				remSpan := strom.CreateText("span", formatSeconds(remaining), atoms.FontWeightBold)
				durationItems.Append(remSpan)
				durSpan.SetStyle(map[string]string{"color": colors.Gray})
			} else {
				if !ended {
					durSpan.AddClass("fw-b")
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
		SetStyle(map[string]string{
			"row-gap": sizes.XSmall,
		})
	titlePropertiesStack.Append(propertiesStack)

	for _, p := range propertiesOrder {
		v := vsp[p]
		if v == "" {
			continue
		}

		ptv := propertyTitles[p] + ": " + v

		propertyRow := strom.CreateText("span", ptv).
			SetStyle(map[string]string{
				"color":     colors.Gray,
				"font-size": font_sizes.XSmall,
			})

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

	tileContainer := strom.Create("a", atoms.DisplayFlex, atoms.FlexDirColumn, atoms.BorderRadiusSmall).
		SetAttribute("href", path.Join("/channel", channelId)).
		SetStyle(map[string]string{
			"flow-shrink": "0",
			"padding":     calc.Mult(sizes.Small, 1.5),
			"row-gap":     sizes.XXSmall,
			"background":  colors.Highlight,
			"width":       "max-content",
		})

	var title string
	if tp, ok := rdx.GetLastVal(data.ChannelTitleProperty, channelId); ok && tp != "" {
		title = tp
	}

	tileContainer.Append(strom.CreateText("span", title, atoms.FontWeightBold))

	var newSubtitle string
	cnev := yeti.ChannelNotEndedVideos(channelId, math.MaxInt, rdx)
	if len(cnev) > 0 {
		switch len(cnev) {
		case 1:
			newSubtitle = "1 new video"
		default:
			newSubtitle = strconv.Itoa(len(cnev)) + " new videos"
		}
	} else {
		newSubtitle = "No new videos"
	}

	tileContainer.Append(strom.CreateText("span", newSubtitle).
		SetStyle(map[string]string{
			"font-size": font_sizes.XSmall,
			"color":     colors.Gray,
		}))

	return tileContainer
}
