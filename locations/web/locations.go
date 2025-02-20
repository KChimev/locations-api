package main

import (
	"net/http"

	"github.com/kchimev/locations-api/locations/internal/constants"
	"github.com/kchimev/locations-api/locations/internal/models"
)

type LocationsService struct {
	locEnt models.LocationInterface
	poiEnt models.POIInterface
}

type RequestPayload struct {
	latitude  string
	longitude string
}

type ResponsePayload struct {
	Locations []models.Location
	POIs      []models.POI
}

func clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (a *LocationsService) fetchLocationData(lat, lon float64) (*ResponsePayload, error) {
	locChan := make(chan []models.Location)
	poiChan := make(chan []models.POI)
	errChan := make(chan error, 2)

	go func() {
		locations, err := a.locEnt.Get(lat, lon, constants.DefaultSearchRadius)
		if err != nil {
			errChan <- err
		}
		locChan <- locations
	}()

	go func() {
		pois, err := a.poiEnt.Get(lat, lon, constants.DefaultSearchRadius)
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
		case err := <-errChan:
			return nil, err
		}
	}

	return res, nil
}
