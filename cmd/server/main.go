package main

import (
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	recoverMiddleware "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/mercan/ecommerce/internal/config"
	"github.com/mercan/ecommerce/internal/repositories/rabbitmq"
	"github.com/mercan/ecommerce/internal/routes"
	"github.com/mercan/ecommerce/internal/utils"
	_ "github.com/mercan/ecommerce/swagger"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName:     config.GetServerConfig().AppName,
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	//app.Get("/swagger/*", swagger.HandlerDefault)

	// Recover middleware
	app.Use(recoverMiddleware.New())
	// Logger middleware
	app.Use(logger.New(utils.LoggerConfig()))

	// Close the channel when main function ends
	defer rabbitmq.Close()

	// Setup RabbitMQ Consumer
	go rabbitmq.ConsumeEmailVerificationQueue()
	//go rabbitmq.ConsumePhoneVerificationQueue()

	// Setup User Routes
	routes.SetupUserRoutes(app)

	if err := app.Listen(":" + config.GetServerConfig().Port); err != nil {
		panic(err)
	}
}
