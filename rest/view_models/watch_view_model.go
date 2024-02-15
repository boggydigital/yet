package view_models

import (
	"fmt"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/pasu"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yet/yeti"
	"github.com/boggydigital/yt_urls"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"
)

type WatchViewModel struct {
	VideoId              string
	VideoUrl             string
	CurrentTime          string
	LastEndedTime        string
	VideoPoster          string
	Server               string
	VideoTitle           string
	ChannelId            string
	ChannelTitle         string
	VideoDescription     string
	VideoPropertiesOrder []string
	VideoProperties      map[string]string
	PlaylistViewModel    *PlaylistViewModel
}

var propertyTitles = map[string]string{
	data.VideoViewCountProperty:            "Views",
	data.VideoKeywordsProperty:             "Keywords",
	data.VideoCategoryProperty:             "Category",
	data.VideoUploadDateProperty:           " Uploaded",
	data.VideoPublishDateProperty:          "Published",
	data.VideoDownloadedDateProperty:       "Downloaded",
	data.VideoDurationProperty:             "Duration",
	data.VideoEndedProperty:                "Last Ended",
	data.VideoSkippedProperty:              "Skipped",
	data.VideosWatchlistProperty:           "In Watchlist",
	data.VideosDownloadQueueProperty:       "In Download Queue",
	data.VideoForcedDownloadProperty:       "Forced Download",
	data.VideoSingleFormatDownloadProperty: "Single Format Download",
}

func GetWatchViewModel(videoId, currentTime string, rdx kvas.ReadableRedux) (*WatchViewModel, error) {

	videoUrl, videoTitle, videoDescription := "", "", ""
	//var videoCaptionTracks []yt_urls.CaptionTrack
	playbackType := "streaming"

	if strings.HasSuffix(videoId, yt_urls.DefaultVideoExt) {
		// TODO: check if that local file exists first
		playbackType = "local"
		videoUrl = "/video?file=" + videoId
		//videoTitle = videoId
	}

	videoPoster := fmt.Sprintf("/poster?v=%s&q=%s", videoId, yt_urls.ThumbnailQualityMaxRes)

	// if current time is not specified with a query parameter - read it from metadata
	if currentTime == "" {
		if ct, ok := rdx.GetLastVal(data.VideoProgressProperty, videoId); ok {
			currentTime = ct
		}
	}

	if title, ok := rdx.GetLastVal(data.VideoTitleProperty, videoId); ok && title != "" {
		if channel, ok := rdx.GetLastVal(data.VideoOwnerChannelNameProperty, videoId); ok && channel != "" {
			localVideoFilename := yeti.ChannelTitleVideoIdFilename(channel, title, videoId)
			if absVideosDir, err := pasu.GetAbsDir(paths.Videos); err == nil {
				absLocalVideoFilename := filepath.Join(absVideosDir, localVideoFilename)
				if _, err := os.Stat(absLocalVideoFilename); err == nil {
					playbackType = "local"
					videoUrl = "/video?file=" + url.QueryEscape(localVideoFilename)
					videoTitle = title
					videoDescription, _ = rdx.GetLastVal(data.VideoShortDescriptionProperty, videoId)

					//if vct, err := getLocalCaptionTracks(videoId, rdx); err == nil {
					//	videoCaptionTracks = vct
					//}
				}
			}
		}
	}

	if videoUrl == "" || videoTitle == "" {
		videoPage, err := yeti.GetVideoPage(videoId)
		if err != nil {
			return nil, err
		}

		metadataDir, err := pasu.GetAbsDir(paths.Metadata)
		if err != nil {
			return nil, err
		}

		mdRdx, err := kvas.NewReduxWriter(metadataDir, data.AllProperties()...)
		if err != nil {
			return nil, err
		}

		for p, values := range yeti.ExtractMetadata(videoPage) {
			if err := mdRdx.AddValues(p, videoId, values...); err != nil {
				return nil, err
			}
		}

		vu, err := decode(videoId, videoPage.BestFormat().Url, videoPage.PlayerUrl)
		if err != nil {
			return nil, err
		}

		videoUrl = vu.String()
		videoTitle = videoPage.VideoDetails.Title
		videoDescription = videoPage.VideoDetails.ShortDescription
	}

	lastEndedTime := ""
	if et, ok := rdx.GetLastVal(data.VideoEndedProperty, videoId); ok && et != "" {
		lastEndedTime = et
		titlePrefix := "☑️ "
		if rdx.HasKey(data.VideoSkippedProperty, videoId) {
			titlePrefix = "⏭️ "
		}
		videoTitle = titlePrefix + videoTitle
	}

	joinProperties := []string{
		data.VideoKeywordsProperty,
		data.VideoCategoryProperty,
	}
	properties := []string{
		data.VideoViewCountProperty,
		data.VideoKeywordsProperty,
		data.VideoCategoryProperty,
		data.VideoUploadDateProperty,
		data.VideoPublishDateProperty,
		data.VideoDownloadedDateProperty,
		data.VideoDurationProperty,
		data.VideoEndedProperty,
		data.VideoSkippedProperty,
		data.VideosWatchlistProperty,
		data.VideosDownloadQueueProperty,
		data.VideoForcedDownloadProperty,
		data.VideoSingleFormatDownloadProperty,
	}

	videoProperties := make(map[string]string)
	titles := make([]string, 0, len(properties))

	for _, p := range properties {
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

	playlistId := ""
	if playlistIds := rdx.MatchAsset(data.PlaylistVideosProperty, []string{videoId}, nil); len(playlistIds) > 0 {
		playlistId = playlistIds[0]
	}

	channelId, channelTitle := "", ""
	if ci, ok := rdx.GetLastVal(data.VideoExternalChannelIdProperty, videoId); ok && ci != "" {
		channelId = ci
		if ct, ok := rdx.GetLastVal(data.VideoOwnerChannelNameProperty, videoId); ok && ct != "" {
			channelTitle = ct
		}
	}

	return &WatchViewModel{
		VideoId:              videoId,
		VideoUrl:             videoUrl,
		VideoPoster:          videoPoster,
		Server:               playbackType,
		VideoTitle:           videoTitle,
		VideoDescription:     videoDescription,
		VideoPropertiesOrder: titles,
		VideoProperties:      videoProperties,
		CurrentTime:          currentTime,
		LastEndedTime:        lastEndedTime,
		ChannelId:            channelId,
		ChannelTitle:         channelTitle,
		PlaylistViewModel:    GetPlaylistViewModel(playlistId, rdx),
	}, nil
}

func decode(videoId, urlStr, playerUrl string) (*url.URL, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	np := q.Get("n")
	if dnp, err := yeti.DecodeParam(http.DefaultClient, videoId, np, playerUrl); err != nil {
		return nil, err
	} else {
		q.Set("n", dnp)
		u.RawQuery = q.Encode()
		return u, nil
	}
}

func getLocalCaptionTracks(videoId string, rdx kvas.ReadableRedux) ([]yt_urls.CaptionTrack, error) {

	if err := rdx.MustHave(
		data.VideoCaptionsNamesProperty,
		data.VideoCaptionsKindsProperty,
		data.VideoCaptionsLanguagesProperty); err != nil {
		return nil, err
	}

	captionsNames, _ := rdx.GetAllValues(data.VideoCaptionsNamesProperty, videoId)
	captionsKinds, _ := rdx.GetAllValues(data.VideoCaptionsKindsProperty, videoId)
	captionsLanguages, _ := rdx.GetAllValues(data.VideoCaptionsLanguagesProperty, videoId)

	cts := make([]yt_urls.CaptionTrack, 0, len(captionsLanguages))
	for i := 0; i < len(captionsLanguages); i++ {

		cn, ck, cl := "", "", captionsLanguages[i]
		if len(captionsNames) >= i {
			cn = captionsNames[i]
		}
		if len(captionsKinds) >= i {
			ck = captionsKinds[i]
		}

		ct := yt_urls.CaptionTrack{
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
	case data.VideoDownloadedDateProperty:
		fallthrough
	case data.VideoEndedProperty:
		if dt, err := time.Parse(time.RFC3339, value); err == nil {
			return dt.Format(time.RFC1123)
		}
	default:
		return value
	}
	return value
}
