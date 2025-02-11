package handlers

import (
	"cloud-solutions-api/authentication"
	"cloud-solutions-api/document"
	"cloud-solutions-api/models"
	"context"
	"database/sql"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func (hc *HandlerContext) CreateDocument(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid request payload",
		})
	}

	account, err := authentication.GetCurrentAccount(hc.Queryer, c)
	if err != nil {
		return c.JSON(
			http.StatusUnauthorized, echo.Map{"error": "Unauthorized"},
		)
	}

	path, err := document.SaveDocumentFile(file)
	if err != nil {
		return err
	}

	text := "" // TODO: read TEXT
	newDocument, err := hc.Queryer.CreateDocument(
		context.Background(),
		models.CreateDocumentParams{
			Name:      file.Filename,
			Text:      sql.NullString{String: text, Valid: true},
			FilePath:  sql.NullString{String: path, Valid: true},
			Embedding: nil,
			AccountID: account.ID,
		},
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, newDocument)
}

func (hc *HandlerContext) DeleteDocumentByID(c echo.Context) error {
	documentIDString := c.Param("documentID")

	documentID, err := strconv.Atoi(documentIDString)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid document ID",
		})
	}

	account, err := authentication.GetCurrentAccount(hc.Queryer, c)
	if err != nil {
		return c.JSON(
			http.StatusUnauthorized, echo.Map{"error": "Unauthorized"},
		)
	}

	owned, err := hc.Queryer.AccountOwnsDocument(
		context.Background(),
		models.AccountOwnsDocumentParams{AccountID: account.ID, ID: int32(documentID)},
	)
	if err != nil {
		return err
	}
	if !owned {
		return c.JSON(http.StatusForbidden, echo.Map{
			"error": "Forbidden: You do not own this document",
		})
	}

	err = hc.Queryer.DeleteDocument(
		context.Background(),
		int32(documentID),
	)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{})
}

func RegisterDocumentRoutes(e *echo.Echo, hc *HandlerContext) {
	restricted := echojwt.JWT(hc.Secret)
	documentGroup := e.Group("/documents")
	documentGroup.POST("", hc.CreateDocument, restricted)
	documentGroup.DELETE("/:documentID", hc.DeleteDocumentByID, restricted)
}
