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

		"/watch": BrGzip(GetOnly(Log(http.HandlerFunc(GetWatch)))),
		"/list":  BrGzip(GetOnly(Log(http.HandlerFunc(GetList)))),
		"/new":   BrGzip(GetOnly(Log(http.HandlerFunc(GetNew)))),

		"/progress": PostOnly(Log(http.HandlerFunc(PostProgress))),
		"/ended":    PostOnly(Log(http.HandlerFunc(PostEnded))),

		"/playlist": BrGzip(GetOnly(Log(http.HandlerFunc(GetPlaylist)))),
		"/refresh":  BrGzip(GetOnly(Log(http.HandlerFunc(GetRefresh)))),

		"/": GetOnly(Log(http.RedirectHandler("/list", http.StatusPermanentRedirect))),
	}

	for p, h := range patternHandlers {
		http.HandleFunc(p, h.ServeHTTP)
	}
}
