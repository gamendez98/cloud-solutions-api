package main

import (
	"cloud-solutions-api/config"
	"cloud-solutions-api/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq" // Importing the driver anonymously
)

func main() {
	// Create a new Echo instance
	e := echo.New()
	e.Use(middleware.CORS())

	configuration := config.GetConfig()

	handlerContext := handlers.NewHandlerContext(*configuration)

	// Middleware
	e.Use(middleware.Logger())  // Logs all HTTP requests
	e.Use(middleware.Recover()) // Recovers from panics

	// Routes
	handlers.RegisterAccountRoutes(e, handlerContext)
	handlers.RegisterDocumentRoutes(e, handlerContext)
	handlers.RegisterChatRoutes(e, handlerContext)

	// Start server
	port := ":80"
	if configuration.Development {
		port = ":8080"
	}
	e.Logger.Info("Server is running on port " + port)
	e.Logger.Fatal(e.Start(port))
}
