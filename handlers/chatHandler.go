package handlers

import (
	"cloud-solutions-api/authentication"
	"cloud-solutions-api/models"
	"context"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/sqlc-dev/pqtype"
	"net/http"
	"strconv"
)

func (hc *HandlerContext) CreateEmptyChat(c echo.Context) error {
	account, err := authentication.GetCurrentAccount(hc.Queryer, c)
	if err != nil {
		return c.JSON(
			http.StatusUnauthorized, echo.Map{"error": "Unauthorized"},
		)
	}

	chat, err := hc.Queryer.CreateChat(context.Background(),
		models.CreateChatParams{
			AccountID: account.ID,
			Messages: pqtype.NullRawMessage{
				RawMessage: []byte("[]"),
				Valid:      true,
			},
		},
	)

	return c.JSON(http.StatusCreated, chat)
}

func (hc *HandlerContext) GetChatByID(c echo.Context) error {
	chatIDString := c.Param("chatID")
	chatID, err := strconv.Atoi(chatIDString)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid chat ID",
		})
	}

	account, err := authentication.GetCurrentAccount(hc.Queryer, c)

	if err != nil {
		return c.JSON(
			http.StatusUnauthorized, echo.Map{"error": "Unauthorized"},
		)
	}

	owns, err := hc.Queryer.AccountOwnsChat(context.Background(), models.AccountOwnsChatParams{
		AccountID: account.ID,
		ID:        int32(chatID),
	})

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": err})
	}

	if !owns {
		return c.JSON(http.StatusForbidden, echo.Map{
			"error": "Forbidden: You do not own this chat",
		})
	}

	chat, err := hc.Queryer.GetChatByID(context.Background(), int32(chatID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": err})
	}

	return c.JSON(http.StatusOK, chat)
}

func (hc *HandlerContext) DeleteChatByID(c echo.Context) error {
	chatIDString := c.Param("chatID")
	chatID, err := strconv.Atoi(chatIDString)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid chat ID",
		})
	}

	account, err := authentication.GetCurrentAccount(hc.Queryer, c)

	if err != nil {
		return c.JSON(
			http.StatusUnauthorized, echo.Map{"error": "Unauthorized"},
		)
	}

	owns, err := hc.Queryer.AccountOwnsChat(context.Background(), models.AccountOwnsChatParams{
		AccountID: account.ID,
		ID:        int32(chatID),
	})

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": err})
	}

	if !owns {
		return c.JSON(http.StatusForbidden, echo.Map{
			"error": "Forbidden: You do not own this chat",
		})
	}

	err = hc.Queryer.DeleteChat(context.Background(), int32(chatID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": err})
	}

	return c.JSON(http.StatusOK, echo.Map{})
}

func RegisterChatRoutes(e *echo.Echo, hc *HandlerContext) {
	restricted := echojwt.JWT(hc.Secret)
	chatGroup := e.Group("/chats")
	chatGroup.POST("", hc.CreateEmptyChat, restricted)
	chatGroup.GET("/:chatID", hc.GetChatByID, restricted)
	chatGroup.DELETE("/:chatID", hc.DeleteChatByID, restricted)
}
