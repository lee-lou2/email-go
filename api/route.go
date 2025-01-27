package api

import "github.com/gofiber/fiber/v3"

// setV1Routes V1 Routes
func setV1Routes(app *fiber.App) {
	app.Get("/v1/plans/:planId", getResultCountHandler)
	app.Get("/v1/plans/counts/sent", getSentCountHandler)
	app.Get("/v1/events/open", createOpenEventHandler)
	app.Post("/v1/events/result", createResultEventHandler)
}
