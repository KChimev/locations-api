package main

import (
	"net/http"
)

func (a *app) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/generate-token", a.generate)
	mux.HandleFunc("/verify-token", a.verify)

	return mux
}
