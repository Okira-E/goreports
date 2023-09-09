package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// SwaggerRouter sets up the swagger routes
func SwaggerRouter(app *fiber.App) {
	port := "3200"

	app.Get("/docs/swagger.json", func(ctx *fiber.Ctx) error {
		// Only allow swagger.json to be retrieved from localhost
		if ctx.IP() != "127.0.0.1" && ctx.IP() != "::1" {
			return ctx.Status(fiber.StatusForbidden).SendString("Forbidden")
		}

		return ctx.Status(fiber.StatusOK).SendFile("./docs/swagger.json")
	})

	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL: "http://localhost:" + port + "/docs/swagger.json",
	}))
}
