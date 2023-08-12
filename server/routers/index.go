package routers

import (
	"github.com/Okira-E/goreports/datasource"
	"github.com/gofiber/fiber/v2"
)

var InternalDb *datasource.DataSource
var ExternalDb *datasource.DataSource

func GlobalRouter(app *fiber.App) {
	ReportsRouter(app)
}
