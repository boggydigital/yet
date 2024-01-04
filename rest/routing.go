package rest

import (
	"github.com/boggydigital/middleware"
	"github.com/boggydigital/nod"
	"net/http"
)

var (
	GetOnly  = middleware.GetMethodOnly
	PostOnly = middleware.PostMethodOnly
	Log      = nod.RequestLog
	BrGzip   = middleware.BrGzip
)

func HandleFuncs() {

	patternHandlers := map[string]http.Handler{
		"/video":    GetOnly(Log(http.HandlerFunc(GetVideo))),
		"/poster":   GetOnly(Log(http.HandlerFunc(GetPoster))),
		"/captions": GetOnly(Log(http.HandlerFunc(GetCaptions))),

		"/watch":        BrGzip(GetOnly(Log(http.HandlerFunc(GetWatch)))),
		"/manage_video": BrGzip(GetOnly(http.HandlerFunc(GetManageVideo))),
		"/update_video": BrGzip(GetOnly(http.HandlerFunc(GetUpdateVideo))),

		"/list": BrGzip(GetOnly(Log(http.HandlerFunc(GetList)))),

		"/paste": BrGzip(GetOnly(Log(http.HandlerFunc(GetPaste)))),

		"/history": BrGzip(GetOnly(Log(http.HandlerFunc(GetHistory)))),

		"/search":  BrGzip(GetOnly(Log(http.HandlerFunc(GetSearch)))),
		"/results": BrGzip(GetOnly(Log(http.HandlerFunc(GetResults)))),

		"/progress": PostOnly(Log(http.HandlerFunc(PostProgress))),
		"/ended":    PostOnly(Log(http.HandlerFunc(PostEnded))),

		"/playlist":         BrGzip(GetOnly(Log(http.HandlerFunc(GetPlaylist)))),
		"/refresh_playlist": BrGzip(GetOnly(Log(http.HandlerFunc(GetRefreshPlaylist)))),
		"/manage_playlist":  BrGzip(GetOnly(Log(http.HandlerFunc(GetManagePlaylist)))),
		"/update_playlist":  BrGzip(GetOnly(Log(http.HandlerFunc(GetUpdatePlaylist)))),

		"/channel":         BrGzip(GetOnly(Log(http.HandlerFunc(GetChannel)))),
		"/refresh_channel": BrGzip(GetOnly(Log(http.HandlerFunc(GetRefreshChannel)))),

		"/": GetOnly(Log(http.RedirectHandler("/list", http.StatusPermanentRedirect))),
	}

	for p, h := range patternHandlers {
		http.HandleFunc(p, h.ServeHTTP)
	}
}
