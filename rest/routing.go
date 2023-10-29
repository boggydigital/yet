package rest

import "net/http"

func HandleFuncs() {

	patternHandlers := map[string]http.Handler{
		"/watch": http.HandlerFunc(GetWatch),
		"/":      http.RedirectHandler("/watch", http.StatusPermanentRedirect),
	}

	for p, h := range patternHandlers {
		http.HandleFunc(p, h.ServeHTTP)
	}
}
