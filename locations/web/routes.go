package main

import (
	"net/http"
)

func (a *Application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/get-location", a.handle)

	return a.checkAuthorization((a.logRequest(mux)))
}
