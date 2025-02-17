package main

import (
	"net/http"

	"github.com/go-playground/validator/v10"
)

func (a *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := &payload{latitude: r.URL.Query().Get("lat"), longitude: r.URL.Query().Get("lon")}
		validator := validator.New(validator.WithRequiredStructEnabled())

		err := validator.Struct(payload)
		if err != nil {
			a.clientError(w, http.StatusBadRequest)
			return
		}

		a.infoLog.Printf("Requested location = Latitude: %s, Longitude: %s", payload.latitude, payload.longitude)

		next.ServeHTTP(w, r)
	})
}

func (a *application) checkAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("Authorization")
		if authToken == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		valid, err := a.validateToken(authToken)
		if err != nil || !valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
