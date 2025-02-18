package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/kchimev/locations-api/locations/internal/constants"
	"github.com/kchimev/locations-api/locations/internal/models"
	"github.com/streadway/amqp"
)

type payload struct {
	latitude  string
	longitude string
}

type application struct {
	errLog    *log.Logger
	infoLog   *log.Logger
	locations *models.LocationEntity
	pois      *models.POIEntity
	mqChan    *amqp.Channel
	mqCon     *amqp.Connection
}

func CreateApplication(errLog *log.Logger, infoLog *log.Logger) (*application, error) {
	db, err := GetDBPool()
	if err != nil {
		return nil, err
	}

	poidb, err := GetPOIDbPool()
	if err != nil {
		return nil, err
	}

	mqcon, mqChan, err := setupRabbit()
	if err != nil {
		return nil, err
	}

	app := &application{
		errLog:    errLog,
		infoLog:   infoLog,
		locations: &models.LocationEntity{DB: db},
		pois:      &models.POIEntity{DB: poidb},
		mqChan:    mqChan,
		mqCon:     mqcon,
	}

	return app, nil
}

func main() {
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stderr, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)

	app, err := CreateApplication(errLog, infoLog)
	if err != nil {
		errLog.Fatal(err)
	}
	defer app.mqChan.Close()
	defer app.mqCon.Close()
	defer app.locations.DB.Close()
	defer app.pois.DB.Close()

	server := &http.Server{
		Addr:     constants.DefaultApplicationPort,
		Handler:  app.routes(),
		ErrorLog: errLog,
	}

	go func() {
		log.Println("Server starting...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
}
