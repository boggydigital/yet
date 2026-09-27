package rest

import (
	"net/http"

	"github.com/boggydigital/nod"
)

var (
	Log = nod.RequestLog
)

func HandleFuncs() {

	patternHandlers := map[string]http.Handler{
		"GET /video":  Log(http.HandlerFunc(GetVideo)),
		"GET /poster": Log(http.HandlerFunc(GetPoster)),

		"GET /watch/{videoId}":        Log(http.HandlerFunc(GetWatch)),
		"GET /manage_video/{videoId}": Log(http.HandlerFunc(GetManageVideo)),
		"GET /video_error/{videoId}":  Log(http.HandlerFunc(GetVideoError)),

		"GET /list": Log(http.HandlerFunc(GetList)),

		"GET /paste":  Log(http.HandlerFunc(GetPaste)),
		"POST /paste": Log(http.HandlerFunc(PostPaste)),

		"GET /history": Log(http.HandlerFunc(GetHistory)),

		"GET /search":  Log(http.HandlerFunc(GetSearch)),
		"GET /results": Log(http.HandlerFunc(GetResults)),

		"GET /playlist/{playlistId}":        Log(http.HandlerFunc(GetPlaylist)),
		"GET /manage_playlist/{playlistId}": Log(http.HandlerFunc(GetManagePlaylist)),
		"GET /update_playlist/{playlistId}": Log(http.HandlerFunc(GetUpdatePlaylist)),

		"GET /channel/{channelId}":        Log(http.HandlerFunc(GetChannel)),
		"GET /manage_channel/{channelId}": Log(http.HandlerFunc(GetManageChannel)),
		"GET /update_channel/{channelId}": Log(http.HandlerFunc(GetUpdateChannel)),

		"POST /progress/{videoId}/{currentTime}": Log(http.HandlerFunc(PostProgress)),
		"POST /end/{videoId}/{reason}":           Log(http.HandlerFunc(PostEnd)),
		"POST /queue_download/{videoId}":         Log(http.HandlerFunc(PostQueueDownload)),
		"POST /download_video/{videoId}":         Log(http.HandlerFunc(PostDownloadVideo)),
		"POST /update_video/{videoId}":           Log(http.HandlerFunc(PostUpdateVideo)),
		"POST /refresh_video/{videoId}":          Log(http.HandlerFunc(PostRefreshVideo)),
		"POST /refresh_playlist/{playlistId}":    Log(http.HandlerFunc(PostRefreshPlaylist)),
		"POST /refresh_channel/{channelId}":      Log(http.HandlerFunc(PostRefreshChannel)),

		"GET /": Log(http.RedirectHandler("/list", http.StatusPermanentRedirect)),
	}

	for p, h := range patternHandlers {
		http.HandleFunc(p, h.ServeHTTP)
	}
}
