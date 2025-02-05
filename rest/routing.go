package rest

import (
	"github.com/boggydigital/nod"
	"net/http"
)

var (
	Log = nod.RequestLog
)

func HandleFuncs() {

	patternHandlers := map[string]http.Handler{
		"GET /video":    Log(http.HandlerFunc(GetVideo)),
		"GET /poster":   Log(http.HandlerFunc(GetPoster)),
		"GET /captions": Log(http.HandlerFunc(GetCaptions)),

		"GET /watch":        Log(http.HandlerFunc(GetWatch)),
		"GET /manage_video": http.HandlerFunc(GetManageVideo),
		"GET /update_video": http.HandlerFunc(GetUpdateVideo),
		"GET /video_error":  http.HandlerFunc(GetVideoError),

		"GET /list": Log(http.HandlerFunc(GetList)),

		"GET /paste": Log(http.HandlerFunc(GetPaste)),

		"GET /history": Log(http.HandlerFunc(GetHistory)),

		"GET /search":  Log(http.HandlerFunc(GetSearch)),
		"GET /results": Log(http.HandlerFunc(GetResults)),

		"POST /progress":       Log(http.HandlerFunc(PostProgress)),
		"POST /ended":          Log(http.HandlerFunc(PostEnded)),
		"POST /queue_download": Log(http.HandlerFunc(PostQueueDownload)),

		"GET /playlist":         Log(http.HandlerFunc(GetPlaylist)),
		"GET /refresh_playlist": Log(http.HandlerFunc(GetRefreshPlaylist)),
		"GET /manage_playlist":  Log(http.HandlerFunc(GetManagePlaylist)),
		"GET /update_playlist":  Log(http.HandlerFunc(GetUpdatePlaylist)),

		"GET /channel":                   Log(http.HandlerFunc(GetChannel)),
		"GET /channel_playlists":         Log(http.HandlerFunc(GetChannelPlaylists)),
		"GET /refresh_channel_videos":    Log(http.HandlerFunc(GetRefreshChannelVideos)),
		"GET /refresh_channel_playlists": Log(http.HandlerFunc(GetRefreshChannelPlaylists)),
		"GET /manage_channel":            http.HandlerFunc(GetManageChannel),
		"GET /update_channel":            http.HandlerFunc(GetUpdateChannel),

		"GET /compton": Log(http.HandlerFunc(GetCompton)),

		"GET /": Log(http.RedirectHandler("/list", http.StatusPermanentRedirect)),
	}

	for p, h := range patternHandlers {
		http.HandleFunc(p, h.ServeHTTP)
	}
}
