package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Create a new Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())  // Logs all HTTP requests
	e.Use(middleware.Recover()) // Recovers from panics

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/users/:id", getUser)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

// Handler
func getUser(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusOK, map[string]string{
		"id":   id,
		"name": "John Doe",
	})
}
