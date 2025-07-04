package view_models

import (
	"fmt"
	"github.com/boggydigital/pathways"
	"github.com/boggydigital/redux"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
	"github.com/boggydigital/yet_urls/youtube_urls"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
)

type WatchViewModel struct {
	VideoId              string
	VideoUrl             string
	CurrentTime          string
	EndedTime            string
	EndedReason          data.VideoEndedReason
	VideoPoster          string
	LocalPlayback        bool
	CurrentTimeSeconds   string
	DurationSeconds      string
	VideoTitle           string
	VideoDescription     string
	VideoPropertiesOrder []string
	VideoProperties      map[string]string
	ChannelViewModel     *ChannelViewModel
	PlaylistViewModel    *PlaylistViewModel
}

var propertyTitles = map[string]string{
	data.VideoViewCountProperty:         "Views",
	data.VideoKeywordsProperty:          "Keywords",
	data.VideoCategoryProperty:          "Category",
	data.VideoUploadDateProperty:        "Uploaded",
	data.VideoPublishDateProperty:       "Published",
	data.VideoDurationProperty:          "Duration",
	data.VideoEndedDateProperty:         "Ended Date",
	data.VideoEndedReasonProperty:       "Ended Reason",
	data.VideoDownloadQueuedProperty:    "Download Queued",
	data.VideoDownloadStartedProperty:   "Download Started",
	data.VideoDownloadCompletedProperty: "Download Completed",
	data.VideoDownloadCleanedUpProperty: "Download Cleaned Up",
	data.VideoForcedDownloadProperty:    "Forced Download",
}

func GetWatchViewModel(videoId, currentTime string, rdx redux.Writeable) (*WatchViewModel, error) {

	videoUrl, videoTitle, videoDescription := "", "", ""
	//var videoCaptionTracks []youtube_urls.CaptionTrack
	localPlayback := false

	if strings.HasSuffix(videoId, youtube_urls.DefaultVideoExt) {
		// TODO: check if that local file exists first
		localPlayback = true
		videoUrl = "/video?file=" + videoId
		//videoTitle = videoId
	}

	videoPoster := fmt.Sprintf("/poster?v=%s&q=%s", videoId, youtube_urls.ThumbnailQualityMaxRes)

	var rem, dur int64

	if durs, sure := rdx.GetLastVal(data.VideoDurationProperty, videoId); sure && durs != "" {
		if duri, err := strconv.ParseInt(durs, 10, 64); err == nil {
			dur = duri
		}
	}

	var ct int64
	if currentTime == "" {

		var currentTimeStr string
		if vpct, ok := data.VideosProgress[videoId]; ok && len(vpct) > 0 {
			currentTimeStr = vpct[0]
		} else if cts, sure := rdx.GetLastVal(data.VideoProgressProperty, videoId); sure && cts != "" {
			currentTimeStr = cts
		}

		if currentTimeStr != "" {
			if cti, err := strconv.ParseInt(currentTimeStr, 10, 64); err == nil {
				ct = cti
				currentTime = strconv.FormatInt(ct, 10)
			}
		}

	}
	rem = dur - ct

	if title, ok := rdx.GetLastVal(data.VideoTitleProperty, videoId); ok && title != "" {
		if channel, ok := rdx.GetLastVal(data.VideoOwnerChannelNameProperty, videoId); ok && channel != "" {
			if absLocalVideoFilename, err := yeti.LocateLocalVideo(videoId); os.IsNotExist(err) {
				// do nothing
			} else if err != nil {
				return nil, err
			} else {
				if _, err := os.Stat(absLocalVideoFilename); err == nil {
					localPlayback = true

					videosDir, err := pathways.GetAbsDir(data.Videos)
					if err != nil {
						return nil, err
					}

					relLocalVideoFilename, err := filepath.Rel(videosDir, absLocalVideoFilename)

					videoUrl = "/video?file=" + url.QueryEscape(relLocalVideoFilename)
					videoTitle = title
					videoDescription, _ = rdx.GetLastVal(data.VideoShortDescriptionProperty, videoId)

					//if vct, err := getLocalCaptionTracks(videoId, rdx); err == nil {
					//	videoCaptionTracks = vct
					//}
				}
			}
		}
	}

	var videoPage *youtube_urls.InitialPlayerResponse
	var err error

	if videoTitle == "" {
		videoPage, err = yeti.GetVideoPage(videoId)
		if err != nil {
			return nil, err
		}

		videoMetadata := yeti.ExtractMetadata(videoPage)

		for p, values := range videoMetadata {
			if err := rdx.AddValues(p, videoId, values...); err != nil {
				return nil, err
			}
		}

		// also set channel title, since ChannelViewModel expects it
		if cids := videoMetadata[data.VideoExternalChannelIdProperty]; len(cids) > 0 {
			if cts := videoMetadata[data.VideoOwnerChannelNameProperty]; len(cts) > 0 {
				channelId := cids[0]
				channelTitle := cts[0]
				if err := rdx.AddValues(data.ChannelTitleProperty, channelId, channelTitle); err != nil {
					return nil, err
				}
			}
		}

		videoTitle = videoPage.VideoDetails.Title
		videoDescription = videoPage.VideoDetails.ShortDescription
	}

	lastEndedTime := ""
	if et, ok := rdx.GetLastVal(data.VideoEndedDateProperty, videoId); ok && et != "" {
		lastEndedTime = et
	}

	endedReason := data.DefaultEndedReason
	if er, ok := rdx.GetLastVal(data.VideoEndedReasonProperty, videoId); ok {
		endedReason = data.ParseVideoEndedReason(er)
	}

	joinProperties := []string{
		data.VideoKeywordsProperty,
		data.VideoCategoryProperty,
	}
	videoProperties := make(map[string]string)
	titles := make([]string, 0, len(propertyTitles))

	for p := range propertyTitles {
		title := propertyTitles[p]
		titles = append(titles, title)
		if slices.Contains(joinProperties, p) {
			if values, ok := rdx.GetAllValues(p, videoId); ok && len(values) > 0 {
				videoProperties[title] = strings.Join(values, ", ")
			}
		} else {
			if value, ok := rdx.GetLastVal(p, videoId); ok && value != "" {
				videoProperties[title] = fmtPropertyValue(p, value)
			}
		}
	}

	sort.Strings(titles)

	allPlaylistsWithVideo := rdx.MatchAsset(data.PlaylistVideosProperty, []string{videoId}, nil)
	playlistId := ""
	for pid := range allPlaylistsWithVideo {
		if rdx.HasKey(data.PlaylistAutoRefreshProperty, pid) {
			playlistId = pid
			break
		}
	}

	if playlistId == "" {
		for pid := range allPlaylistsWithVideo {
			playlistId = pid
			break
		}
	}

	channelId := ""
	if ci, ok := rdx.GetLastVal(data.VideoExternalChannelIdProperty, videoId); ok && ci != "" {
		channelId = ci
	}

	return &WatchViewModel{
		VideoId:              videoId,
		VideoUrl:             videoUrl,
		VideoPoster:          videoPoster,
		LocalPlayback:        localPlayback,
		CurrentTimeSeconds:   strconv.FormatInt(dur-rem, 10),
		DurationSeconds:      strconv.FormatInt(dur, 10),
		VideoTitle:           videoTitle,
		VideoDescription:     videoDescription,
		VideoPropertiesOrder: titles,
		VideoProperties:      videoProperties,
		EndedTime:            lastEndedTime,
		EndedReason:          endedReason,
		ChannelViewModel:     GetChannelViewModel(channelId, rdx),
		PlaylistViewModel:    GetPlaylistViewModel(playlistId, rdx),
	}, nil
}

func getLocalCaptionTracks(videoId string, rdx redux.Readable) ([]youtube_urls.CaptionTrack, error) {

	if err := rdx.MustHave(
		data.VideoCaptionsNamesProperty,
		data.VideoCaptionsKindsProperty,
		data.VideoCaptionsLanguagesProperty); err != nil {
		return nil, err
	}

	captionsNames, _ := rdx.GetAllValues(data.VideoCaptionsNamesProperty, videoId)
	captionsKinds, _ := rdx.GetAllValues(data.VideoCaptionsKindsProperty, videoId)
	captionsLanguages, _ := rdx.GetAllValues(data.VideoCaptionsLanguagesProperty, videoId)

	cts := make([]youtube_urls.CaptionTrack, 0, len(captionsLanguages))
	for i := 0; i < len(captionsLanguages); i++ {

		cn, ck, cl := "", "", captionsLanguages[i]
		if len(captionsNames) >= i {
			cn = captionsNames[i]
		}
		if len(captionsKinds) >= i {
			ck = captionsKinds[i]
		}

		ct := youtube_urls.CaptionTrack{
			BaseUrl:      "/captions?v=" + videoId + "&l=" + cl,
			LanguageCode: cl,
			Kind:         ck,
			TrackName:    cn,
		}

		if ct.Kind == "asr" {
			ct.Kind = "subtitles"
		}

		cts = append(cts, ct)
	}

	return cts, nil
}

func fmtPropertyValue(property, value string) string {
	switch property {
	case data.VideoDurationProperty:
		if iv, err := strconv.ParseInt(value, 10, 32); err == nil {
			return formatSeconds(iv)
		}
	case data.VideoUploadDateProperty:
		fallthrough
	case data.VideoPublishDateProperty:
		fallthrough
	case data.VideoDownloadQueuedProperty:
		fallthrough
	case data.VideoDownloadStartedProperty:
		fallthrough
	case data.VideoDownloadCompletedProperty:
		fallthrough
	case data.VideoDownloadCleanedUpProperty:
		fallthrough
	case data.VideoEndedDateProperty:
		if dt, err := time.Parse(time.RFC3339, value); err == nil {
			return dt.Format(time.RFC1123)
		}
	default:
		return value
	}
	return value
}
