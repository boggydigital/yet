package rest

import "net/http"

func HandleFuncs() {

	patternHandlers := map[string]http.Handler{
		"/local_video": http.HandlerFunc(GetLocalVideo),
		"/watch":       http.HandlerFunc(GetWatch),
		"/list":        http.HandlerFunc(GetList),
		"/":            http.RedirectHandler("/list", http.StatusPermanentRedirect),
	}

	for p, h := range patternHandlers {
		http.HandleFunc(p, h.ServeHTTP)
	}
}
