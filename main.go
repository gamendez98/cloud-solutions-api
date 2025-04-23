package main

import (
	"cloud-solutions-api/config"
	"cloud-solutions-api/handlers"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq" // Importing the driver anonymously
	"net/http"
)

func customHTTPErrorHandler(err error, c echo.Context) {
	var he *echo.HTTPError
	code := http.StatusInternalServerError

	// Try to cast to *echo.HTTPError
	if errors.As(err, &he) {
		code = he.Code
	}

	// If it's a 500 or not an HTTP error at all, log stack trace
	if code == http.StatusInternalServerError || !errors.As(err, &he) {
		c.Logger().Errorf("internal error: %v\n", err)
	}

	// Send a JSON response to the client
	_ = c.JSON(code, map[string]any{
		"error": http.StatusText(code),
	})
}

func main() {
	// Create a new Echo instance
	e := echo.New()
	e.Use(middleware.CORS())
	e.HTTPErrorHandler = customHTTPErrorHandler

	configuration := config.GetConfig()

	handlerContext := handlers.NewHandlerContext(*configuration)
	defer func() {
		errs := handlerContext.Close()
		for _, err := range errs {
			e.Logger.Errorf("error closing handler context: %s", err)
		}
	}()

	// Middleware
	e.Use(middleware.Logger())  // Logs all HTTP requests
	e.Use(middleware.Recover()) // Recovers from panics

	// Routes
	handlers.RegisterAccountRoutes(e, handlerContext)
	handlers.RegisterDocumentRoutes(e, handlerContext)
	handlers.RegisterChatRoutes(e, handlerContext)
	e.GET("/health", handlerContext.HealthCheck)

	// Start server
	port := ":80"
	if configuration.Development {
		port = ":8080"
	}
	e.Logger.Info("Server is running on port " + port)
	e.Logger.Fatal(e.Start(port))
}
