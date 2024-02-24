package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/okira-e/goreports/datasource"
)

var InternalDb *datasource.DataSource
var ExternalDb *datasource.DataSource

func GlobalRouter(app *fiber.App) {
	ReportsRouter(app)
	SwaggerRouter(app)
}
