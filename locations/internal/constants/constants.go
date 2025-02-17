package constants

import (
	"fmt"
)

const (
	// Application
	DefaultApplicationPort = ":4321"
	DefaultSearchRadius    = 10000
	// RabbitMQ
	RabbitURL               = "amqp://guest:guest@localhost:5672/"
	RabbitExchange          = "application_exchange"
	AuthenticationExchange  = "authentication_exchange"
	CurrentServiceQueueName = "locations-queue"
	CurrentServiceQueueKey  = "locations.get"
	AuthenticateQueueName   = "authentication-queue"
	AuthenticateQueueKey    = "authentication.authenticate"
	GenerateQueueName       = "generate-queue"
	GenerateQueueKey        = "authentication.generate"
	AuthenticationTimeout   = 3
	// Database
	DBHost     = "localhost"
	DBName     = "osm_bulgaria"
	DBUser     = "postgres"
	DBPassword = "admin123"
)

var (
	DBConnectionString = fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable",
		DBHost, DBName, DBUser, DBPassword)
)
