package handlers

import (
	"cloud-solutions-api/authentication"
	"cloud-solutions-api/document"
	"cloud-solutions-api/models"
	"cloud-solutions-api/pubSubPublisher"
	"context"
	"database/sql"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

// UserOwnsDocumentMiddleware is a middleware function to check if the currently
// authenticated user owns the document specified in the request.
func (hc *HandlerContext) UserOwnsDocumentMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		documentIDString := c.Param("documentID")

		// Validate document ID format
		documentID, err := strconv.Atoi(documentIDString)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid document ID")
		}

		// Retrieve the current account
		account, err := authentication.GetCurrentAccount(hc.Queryer, c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}

		// Check if the user owns the document
		owned, err := hc.Queryer.AccountOwnsDocument(
			context.Background(),
			models.AccountOwnsDocumentParams{AccountID: account.ID, ID: int32(documentID)},
		)
		if err != nil {
			return err
		}

		// Deny access if the user doesn't own the document
		if !owned {
			return echo.NewHTTPError(http.StatusForbidden, "Forbidden: You do not own this document")
		}

		// If the user owns the document, proceed to the next handler
		return next(c)
	}
}

func (hc *HandlerContext) CreateDocument(c echo.Context) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	account, err := authentication.GetCurrentAccount(hc.Queryer, c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	path, err := document.SaveDocumentFileInBucket(fileHeader, hc.Bucket)
	if err != nil {
		return err
	}

	text, err := document.ExtractTextFromDocumentFile(fileHeader)
	if err != nil {
		c.Logger().Errorf("error extracting text from document: %s", err)
	}

	newDocument, err := hc.Queryer.CreateDocument(
		context.Background(),
		models.CreateDocumentParams{
			Name:      fileHeader.Filename,
			Text:      sql.NullString{String: text, Valid: true},
			FilePath:  sql.NullString{String: path, Valid: true},
			Embedding: nil,
			AccountID: account.ID,
		},
	)
	if err != nil {
		return err
	}

	err = hc.PuSubPublisher.PublishDocumentIndexingMessage(pubSubPublisher.DocumentIndexingMessage{
		DocumentId:   newDocument.ID,
		DocumentText: newDocument.Text.String,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, newDocument)
}

func (hc *HandlerContext) DeleteDocumentByID(c echo.Context) error {
	documentIDString := c.Param("documentID")

	documentID, err := strconv.Atoi(documentIDString)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid document ID")
	}

	retrievedDocument, err := hc.Queryer.GetDocumentByID(context.Background(), int32(documentID))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid document ID")
	}

	err = hc.Queryer.DeleteDocument(
		context.Background(),
		int32(documentID),
	)

	if err != nil {
		return err
	}

	err = document.DeleteDocumentFileFromBucket(retrievedDocument.FilePath.String, hc.Bucket)
	if err != nil {
		c.Logger().Errorf("error deleting document file from bucket: %s", err)
	}

	return c.JSON(http.StatusOK, echo.Map{})
}

// RegisterDocumentRoutes sets up the routes for document operations, applying JWT authentication for restricted access.
func RegisterDocumentRoutes(e *echo.Echo, hc *HandlerContext) {
	restricted := echojwt.JWT(hc.Secret)
	documentGroup := e.Group("/documents")
	documentGroup.POST("", hc.CreateDocument, restricted)
	documentGroup.DELETE("/:documentID", hc.DeleteDocumentByID, restricted, hc.UserOwnsDocumentMiddleware)
}
