package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/kchimev/locations-api/locations/internal/constants"
	"github.com/kchimev/locations-api/locations/internal/models"
)

type ResponsePayload struct {
	Locations []models.Location
	POIs      []models.POI
}

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

	locChan := make(chan []models.Location)
	poiChan := make(chan []models.POI)
	errChan := make(chan error, 2)

	go func() {
		locations, err := a.locations.Get(lat, lon, constants.DefaultSearchRadius)
		if err != nil {
			errChan <- err
		}
		locChan <- locations
	}()

	go func() {
		pois, err := a.pois.Get(lat, lon, constants.DefaultSearchRadius)
		if err != nil {
			errChan <- err
		}
		poiChan <- pois
	}()

	res := &ResponsePayload{}
	for i := 0; i < 2; i++ {
		select {
		case locs, ok := <-locChan:
			if ok {
				res.Locations = locs
			}
		case pois, ok := <-poiChan:
			if ok {
				res.POIs = pois
			}
		case <-errChan:
			a.clientError(w, http.StatusNotFound)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(*res); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}
