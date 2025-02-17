package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/kchimev/locations-api/locations/internal/constants"
)

func (a *application) getLocation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		a.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	payload := &payload{
		latitude:  r.URL.Query().Get("lat"),
		longitude: r.URL.Query().Get("lon"),
	}
	validator := validator.New(validator.WithRequiredStructEnabled())

	err := validator.Struct(payload)
	if err != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}

	lat, err := strconv.ParseFloat(payload.latitude, 64)
	if err != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}

	lon, err := strconv.ParseFloat(payload.longitude, 64)
	if err != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}

	res, err := a.locations.Get(lat, lon, constants.DefaultSearchRadius)
	if err != nil {
		a.clientError(w, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}
