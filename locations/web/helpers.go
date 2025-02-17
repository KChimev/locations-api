package main

import (
	"net/http"
)

func (*application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
