package server

import (
	"github.com/Okira-E/goreports/datasource"
	"github.com/Okira-E/goreports/server/routes"
	"github.com/Okira-E/goreports/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
)

// StartServer starts a Fiber web server to listen for any requests to GoReports.
func StartServer() {
	var internalDb datasource.DataSource
	var externalDb datasource.DataSource

	// Establish a connection with internal database.
	dataDir, errOpt := utils.GetDataDirBasedOnOS()
	if errOpt.IsSome() {
		log.Fatalf("error while getting the data directory: %v", errOpt.Unwrap())
	}

	internalDb = datasource.NewSqliteDb(dataDir + "/internal.db")

	errOpt = internalDb.Connect()
	if errOpt.IsSome() {
		log.Fatalf("error while connecting to the database: %v", errOpt.Unwrap())
	}

	// Establish a connection with external database.
	config, errOpt := utils.GetConfigData()
	if errOpt.IsSome() {
		log.Fatalf("error while getting the config data: %v", errOpt.Unwrap())
	}

	connStr, errMsgOpt := config.DbConfig.GetConnectionString()
	if errMsgOpt.IsSome() {
		log.Fatalf("error while getting the connection string: %v", errMsgOpt.Unwrap())
	}
	externalDb = datasource.NewExternalDb(connStr)

	errOpt = externalDb.Connect()
	if errOpt.IsSome() {
		log.Fatalf("error while connecting to the external database: %v", errOpt.Unwrap())
	}

	// Close the database connections when the server stops.
	defer func() {
		if errOpt = internalDb.Disconnect(); errOpt.IsSome() {
			log.Fatalf("error while disconnecting from the internal database: %v", errOpt.Unwrap())
		}

		if errOpt = externalDb.Disconnect(); errOpt.IsSome() {
			log.Fatalf("error while disconnecting from the external database: %v", errOpt.Unwrap())
		}
	}()

	// Create a new Fiber instance.
	app := fiber.New()
	// Set up CORS.
	app.Use(cors.New())
	// Set up the databases.
	routes.InternalDb = &internalDb
	routes.ExternalDb = &externalDb
	// Set up the routes.
	routes.GlobalRouter(app)

	// Start the server.
	const port = ":3200"
	err := app.Listen("0.0.0.0" + port)
	if err != nil {
		log.Fatalf("error while starting the server: %v", err)
	}
}
