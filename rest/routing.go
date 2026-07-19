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

		"GET /watch/{videoId}":         Log(http.HandlerFunc(GetWatch)),
		"GET /refresh_video/{videoId}": Log(http.HandlerFunc(GetRefreshVideo)),
		"GET /manage_video":            http.HandlerFunc(GetManageVideo),
		"GET /update_video":            http.HandlerFunc(GetUpdateVideo),
		"GET /video_error":             http.HandlerFunc(GetVideoError),

		"GET /list": Log(http.HandlerFunc(GetList)),

		"GET /paste":       Log(http.HandlerFunc(GetPaste)),
		"GET /paste_video": Log(http.HandlerFunc(GetPasteVideo)),

		"GET /history": Log(http.HandlerFunc(GetHistory)),

		"GET /search":  Log(http.HandlerFunc(GetSearch)),
		"GET /results": Log(http.HandlerFunc(GetResults)),

		"POST /progress/{videoId}/{currentTime}": Log(http.HandlerFunc(PostProgress)),
		"GET /end/{videoId}/{reason}":            Log(http.HandlerFunc(GetEnded)),
		"GET /queue_download/{videoId}":          Log(http.HandlerFunc(GetQueueDownload)),
		"GET /download_video/{videoId}":          Log(http.HandlerFunc(GetDownloadVideo)),

		"GET /playlist/{playlistId}":         Log(http.HandlerFunc(GetPlaylist)),
		"GET /refresh_playlist/{playlistId}": Log(http.HandlerFunc(GetRefreshPlaylist)),
		"GET /manage_playlist/{playlistId}":  Log(http.HandlerFunc(GetManagePlaylist)),
		"GET /update_playlist/{playlistId}":  Log(http.HandlerFunc(GetUpdatePlaylist)),

		"GET /channel/{channelId}":                   Log(http.HandlerFunc(GetChannel)),
		"GET /channel_playlists/{channelId}":         Log(http.HandlerFunc(GetChannelPlaylists)),
		"GET /refresh_channel_videos/{channelId}":    Log(http.HandlerFunc(GetRefreshChannelVideos)),
		"GET /refresh_channel_playlists/{channelId}": Log(http.HandlerFunc(GetRefreshChannelPlaylists)),
		"GET /manage_channel/{channelId}":            http.HandlerFunc(GetManageChannel),
		"GET /update_channel/{channelId}":            http.HandlerFunc(GetUpdateChannel),

		"GET /": Log(http.RedirectHandler("/list", http.StatusPermanentRedirect)),
	}

	for p, h := range patternHandlers {
		http.HandleFunc(p, h.ServeHTTP)
	}
}
