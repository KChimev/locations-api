package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/kchimev/locations-api/authentication/internal/constants"

	"github.com/streadway/amqp"
)

type app struct {
	mqcon  *amqp.Connection
	mqchan *amqp.Channel
}

func main() {
	mqcon, mqchan, err := setupRabbit()
	if err != nil {
		log.Fatalf("Rabbitmq connection failed: %v", err)
	}
	defer mqcon.Close()
	defer mqchan.Close()
	app := &app{mqchan: mqchan, mqcon: mqcon}
	server := &http.Server{
		Handler: app.routes(),
		Addr:    constants.DefaultApplicationPort,
	}

	go app.consumeQueue(constants.AuthenticateQueueName, app.handleAuthenticate)
	go app.consumeQueue(constants.GenerateQueueName, app.handleGenerate)
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
