package rest

import (
	"github.com/boggydigital/middleware"
	"github.com/boggydigital/nod"
	"net/http"
)

var (
	Log    = nod.RequestLog
	BrGzip = middleware.BrGzip
)

func HandleFuncs() {

	patternHandlers := map[string]http.Handler{
		"GET /video":    Log(http.HandlerFunc(GetVideo)),
		"GET /poster":   Log(http.HandlerFunc(GetPoster)),
		"GET /captions": Log(http.HandlerFunc(GetCaptions)),

		"GET /watch":        BrGzip(Log(http.HandlerFunc(GetWatch))),
		"GET /manage_video": BrGzip(http.HandlerFunc(GetManageVideo)),
		"GET /update_video": BrGzip(http.HandlerFunc(GetUpdateVideo)),
		"GET /video_error":  BrGzip(http.HandlerFunc(GetVideoError)),

		"GET /list": BrGzip(Log(http.HandlerFunc(GetList))),

		"GET /paste": BrGzip(Log(http.HandlerFunc(GetPaste))),

		"GET /history": BrGzip(Log(http.HandlerFunc(GetHistory))),

		"GET /search":  BrGzip(Log(http.HandlerFunc(GetSearch))),
		"GET /results": BrGzip(Log(http.HandlerFunc(GetResults))),

		"POST /progress": Log(http.HandlerFunc(PostProgress)),
		"POST /ended":    Log(http.HandlerFunc(PostEnded)),

		"GET /playlist":         BrGzip(Log(http.HandlerFunc(GetPlaylist))),
		"GET /refresh_playlist": BrGzip(Log(http.HandlerFunc(GetRefreshPlaylist))),
		"GET /manage_playlist":  BrGzip(Log(http.HandlerFunc(GetManagePlaylist))),
		"GET /update_playlist":  BrGzip(Log(http.HandlerFunc(GetUpdatePlaylist))),

		"GET /channel":         BrGzip(Log(http.HandlerFunc(GetChannel))),
		"GET /refresh_channel": BrGzip(Log(http.HandlerFunc(GetRefreshChannel))),

		"GET /": Log(http.RedirectHandler("/list", http.StatusPermanentRedirect)),
	}

	for p, h := range patternHandlers {
		http.HandleFunc(p, h.ServeHTTP)
	}
}
