package constants

const (
	// Application
	DefaultApplicationPort = ":1234"
	JWTKey                 = "SomeSecretKey"
	// RabbitMQ
	RabbitURL             = "amqp://guest:guest@localhost:5672/"
	RabbitExchange        = "authentication_exchange"
	AuthenticateQueueName = "authentication-queue"
	AuthenticateQueueKey  = "authentication.authenticate"
	GenerateQueueName     = "generate-queue"
	GenerateQueueKey      = "authentication.generate"
)
