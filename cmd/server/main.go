package main

import (
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	recoverMiddleware "github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/mercan/ecommerce/internal/config"
	"github.com/mercan/ecommerce/internal/helpers"
	"github.com/mercan/ecommerce/internal/repositories/rabbitmq"
	"github.com/mercan/ecommerce/internal/routes"
)

// main is the entry point of the application
func main() {
	// Create a new Fiber app with configuration

	app := fiber.New(fiber.Config{
		AppName:       config.GetServerConfig().AppName,
		ServerHeader:  "",
		CaseSensitive: true,
		JSONEncoder:   json.Marshal,
		JSONDecoder:   json.Unmarshal,
	})

	// Use recover and logger middlewares
	app.Use(recoverMiddleware.New())
	app.Use(logger.New(helpers.LoggerConfig()))

	// Defer closing the RabbitMQ channel when the main function ends
	defer rabbitmq.Close()

	// Setup RabbitMQ Consumers for email and phone verification queues
	emailQueue := rabbitmq.NewEmailQueueManager()
	go emailQueue.ConsumeEmailVerificationQueue()
	phoneQueue := rabbitmq.NewPhoneQueueManager()
	go phoneQueue.ConsumePhoneVerificationQueue()

	// Setup User Routes
	routes.SetupUserRoutes(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// Listen on the configured server port
	if err := app.Listen(":" + config.GetServerConfig().Port); err != nil {
		panic(err)
	}
}
