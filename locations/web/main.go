package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/go-playground/validator/v10"
	"github.com/kchimev/locations-api/locations/internal/constants"
	"github.com/kchimev/locations-api/locations/internal/models"
)

type Application struct {
	errLog    *log.Logger
	infoLog   *log.Logger
	locations *LocationsService
	rabbit    *RabbitConnection
}

func CreateApplication(errLog *log.Logger, infoLog *log.Logger) (*Application, error) {
	locDb, err := GetLocationsDBPool()
	if err != nil {
		return nil, err
	}

	poiDb, err := GetPOIDbPool()
	if err != nil {
		return nil, err
	}

	mqCon, mqChan, err := setupRabbit()
	if err != nil {
		return nil, err
	}

	server := &Application{
		errLog:  errLog,
		infoLog: infoLog,
		locations: &LocationsService{
			locEnt: &models.LocationEntity{locDb},
			poiEnt: &models.POIEntity{poiDb},
		},
		rabbit: &RabbitConnection{
			mqChan: mqChan,
			mqCon:  mqCon,
		},
	}

	return server, nil
}

func (a *Application) handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		clientError(w, http.StatusMethodNotAllowed)
		return
	}

	payload := &RequestPayload{
		latitude:  r.URL.Query().Get("lat"),
		longitude: r.URL.Query().Get("lon"),
	}
	validator := validator.New(validator.WithRequiredStructEnabled())

	err := validator.Struct(payload)
	if err != nil {
		clientError(w, http.StatusBadRequest)
		return
	}

	lat, err := strconv.ParseFloat(payload.latitude, 64)
	if err != nil {
		clientError(w, http.StatusBadRequest)
		return
	}

	lon, err := strconv.ParseFloat(payload.longitude, 64)
	if err != nil {
		clientError(w, http.StatusBadRequest)
		return
	}

	res, err := a.locations.fetchLocationData(lat, lon)
	if err != nil {
		clientError(w, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(*res); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

func main() {
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stderr, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)

	app, err := CreateApplication(errLog, infoLog)
	if err != nil {
		errLog.Fatal(err)
	}
	defer app.rabbit.mqChan.Close()
	defer app.rabbit.mqCon.Close()
	defer app.locations.locEnt.Close()
	defer app.locations.poiEnt.Close()

	server := &http.Server{
		Addr:     constants.DefaultApplicationPort,
		Handler:  app.routes(),
		ErrorLog: errLog,
	}

	go func() {
		app.infoLog.Println("Server starting...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.errLog.Fatalf("HTTP server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
}
