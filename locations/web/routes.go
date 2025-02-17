package main

import (
	"net/http"
)

func (a *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/get-location", a.getLocation)

	return a.checkAuthorization((a.logRequest(mux)))
}
