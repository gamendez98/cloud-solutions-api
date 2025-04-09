package handlers

import (
	"cloud-solutions-api/authentication"
	"cloud-solutions-api/models"
	"context"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
)

func getOffsetLimit(c echo.Context) (int, int) {
	offsetString := c.QueryParam("offset")
	limitString := c.QueryParam("limit")
	offset, err := strconv.Atoi(offsetString)
	if err != nil {
		offset = 0
	}
	limit, err := strconv.Atoi(limitString)
	if err != nil {
		limit = 10
	}
	return offset, limit
}

func (hc *HandlerContext) login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	message401 := "Invalid username or password"

	passwordHash, err := hc.Queryer.GetAccountPasswordHashByUsername(context.Background(), username)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": message401,
		})
	}

	if !authentication.CheckPasswordHash(password, passwordHash.PasswordHash) {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": message401,
		})
	}

	// Generate JWT token
	claims := &authentication.JwtCustomClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // 1-day expiration
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(hc.Secret)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": signedToken,
	})
}

func (hc *HandlerContext) CreateUser(c echo.Context) error {
	var accountCreationParams = struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}{}
	if err := c.Bind(&accountCreationParams); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid request payload",
		})
	}

	password_hash, err := authentication.HashPassword(accountCreationParams.Password)
	if err != nil {
		return err
	}

	account, err := hc.Queryer.CreateAccount(
		context.Background(),
		models.CreateAccountParams{
			Username:     accountCreationParams.Username,
			Email:        accountCreationParams.Email,
			PasswordHash: password_hash,
		},
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, account)
}

func (hc *HandlerContext) GetAccountByID(c echo.Context) error {
	username, err := authentication.GetCurrentUsername(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}
	account, err := hc.Queryer.GetAccountByUsername(context.Background(), username)
	if err != nil {
		return err
	}
	return c.JSON(200, account)
}

func (hc *HandlerContext) GetAccountDocuments(c echo.Context) error {
	account, err := authentication.GetCurrentAccount(hc.Queryer, c)
	offset, limit := getOffsetLimit(c)
	if err != nil {
		return c.JSON(
			http.StatusUnauthorized, echo.Map{"error": "Unauthorized"},
		)
	}
	documents, err := hc.Queryer.GetDocumentsByAccountID(
		context.Background(),
		models.GetDocumentsByAccountIDParams{
			AccountID: account.ID,
			Offset:    int32(offset),
			Limit:     int32(limit),
		},
	)
	if err != nil {
		return err
	}
	return c.JSON(
		http.StatusOK, documents,
	)
}

func (hc *HandlerContext) GetAccountChats(c echo.Context) error {
	account, err := authentication.GetCurrentAccount(hc.Queryer, c)
	offset, limit := getOffsetLimit(c)
	if err != nil {
		return c.JSON(
			http.StatusUnauthorized, echo.Map{"error": "Unauthorized"},
		)
	}

	chats, err := hc.Queryer.GetChatsByAccountID(
		context.Background(),
		models.GetChatsByAccountIDParams{
			AccountID: account.ID,
			Offset:    int32(offset),
			Limit:     int32(limit),
		},
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, chats)
}

// RegisterAccountRoutes registers account-related routes
func RegisterAccountRoutes(e *echo.Echo, hc *HandlerContext) {
	restricted := echojwt.JWT(hc.Secret)
	accountGroup := e.Group("/accounts")
	accountGroup.POST("/login", hc.login)
	accountGroup.POST("", hc.CreateUser)
	accountGroup.GET("", hc.GetAccountByID, restricted)
	accountGroup.GET("/documents", hc.GetAccountDocuments, restricted)
	accountGroup.GET("/chats", hc.GetAccountChats, restricted)
}
