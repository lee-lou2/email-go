package api

import (
	"aws-ses-sender-go/config"
	"fmt"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func Run() {
	app := fiber.New(
		fiber.Config{AppName: "aws-ses-sender-go"},
	)

	// Middleware
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format:     "${time} ${pid} ${locals:requestid} ${status} - ${method} ${path} ${latency}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	}))

	// Routes
	setV1Routes(app)

	log.Fatal(app.Listen(fmt.Sprintf(":%s", config.GetEnv("SERVER_PORT", "3000"))))
}
