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
)

func HandleFuncs() {

	patternHandlers := map[string]http.Handler{
		"/video": GetOnly(Log(http.HandlerFunc(GetVideo))),
		"/watch": GetOnly(Log(http.HandlerFunc(GetWatch))),
		"/list":  GetOnly(Log(http.HandlerFunc(GetList))),

		"/progress": PostOnly(Log(http.HandlerFunc(PostProgress))),
		"/ended":    PostOnly(Log(http.HandlerFunc(PostEnded))),

		"/": GetOnly(Log(http.RedirectHandler("/list", http.StatusPermanentRedirect))),
	}

	for p, h := range patternHandlers {
		http.HandleFunc(p, h.ServeHTTP)
	}
}
