package handlers

import (
	"cloud-solutions-api/models"
	"github.com/labstack/echo/v4"
)

func (hc *HandlerContext) GetAccountByID(c echo.Context) error {
	var account = models.Account{}
	return c.JSON(200, account)
}

// RegisterAccountRoutes registers account-related routes
func RegisterAccountRoutes(e *echo.Echo, hc *HandlerContext) {
	userGroup := e.Group("/accounts")
	userGroup.GET("/:id", hc.GetAccountByID)
}
