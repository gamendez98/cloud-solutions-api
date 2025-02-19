package handlers

import (
	"cloud-solutions-api/authentication"
	"cloud-solutions-api/chat"
	"cloud-solutions-api/models"
	"cloud-solutions-api/rabbitConnection"
	"context"
	"encoding/json"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/sqlc-dev/pqtype"
	"net/http"
	"strconv"
)

func CheckChatOwnership(hc *HandlerContext, c echo.Context, chatID int) bool {
	account, err := authentication.GetCurrentAccount(hc.Queryer, c)

	if err != nil {
		return false
	}

	owns, err := hc.Queryer.AccountOwnsChat(context.Background(), models.AccountOwnsChatParams{
		AccountID: account.ID,
		ID:        int32(chatID),
	})

	if err != nil {
		return false
	}

	return owns
}

func (hc *HandlerContext) CreateEmptyChat(c echo.Context) error {
	account, err := authentication.GetCurrentAccount(hc.Queryer, c)
	if err != nil {
		return c.JSON(
			http.StatusUnauthorized, echo.Map{"error": "Unauthorized"},
		)
	}

	newChat, err := hc.Queryer.CreateChat(context.Background(),
		models.CreateChatParams{
			AccountID: account.ID,
			Messages: pqtype.NullRawMessage{
				RawMessage: []byte("[]"),
				Valid:      true,
			},
		},
	)

	return c.JSON(http.StatusCreated, newChat)
}

func (hc *HandlerContext) GetChatByID(c echo.Context) error {
	chatIDString := c.Param("chatID")
	chatID, err := strconv.Atoi(chatIDString)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid chat ID",
		})
	}

	owns := CheckChatOwnership(hc, c, chatID)

	if !owns {
		return c.JSON(http.StatusForbidden, echo.Map{
			"error": "Forbidden: You do not own this chat",
		})
	}

	retrievedChat, err := hc.Queryer.GetChatByID(context.Background(), int32(chatID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": err})
	}

	return c.JSON(http.StatusOK, retrievedChat)
}

func (hc *HandlerContext) DeleteChatByID(c echo.Context) error {
	chatIDString := c.Param("chatID")
	chatID, err := strconv.Atoi(chatIDString)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid chat ID",
		})
	}

	owns := CheckChatOwnership(hc, c, chatID)

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

func (hc *HandlerContext) CreateChatMessage(c echo.Context) error {
	chatIDString := c.Param("chatID")
	chatID, err := strconv.Atoi(chatIDString)
	newMessageParameters := chat.NewMessageParameters{}

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid chat ID",
		})
	}

	if err = c.Bind(&newMessageParameters); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
	}

	owns := CheckChatOwnership(hc, c, chatID)

	if !owns {
		return c.JSON(http.StatusForbidden, echo.Map{
			"error": "Forbidden: You do not own this chat",
		})
	}

	retrievedChat, err := hc.Queryer.GetChatByID(context.Background(), int32(chatID))

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": err})
	}

	newMessage := chat.NewMessageForChat(
		retrievedChat,
		newMessageParameters,
	)
	newMessageJSON, err := json.Marshal(newMessage)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": err})
	}

	retrievedChat, err = hc.Queryer.AddMessageToChat(context.Background(), models.AddMessageToChatParams{
		Chatid:     retrievedChat.ID,
		Newmessage: newMessageJSON,
	})

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": err})
	}

	err = hc.AIAssistantMessagePublisher.Publish(rabbitConnection.AIAssistantMessage{
		Messages: retrievedChat.GetMessages(),
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error publishing message to RabbitMQ"})
	}

	return c.JSON(http.StatusOK, retrievedChat)
}

func RegisterChatRoutes(e *echo.Echo, hc *HandlerContext) {
	restricted := echojwt.JWT(hc.Secret)
	chatGroup := e.Group("/chats")
	chatGroup.POST("", hc.CreateEmptyChat, restricted)
	chatGroup.GET("/:chatID", hc.GetChatByID, restricted)
	chatGroup.DELETE("/:chatID", hc.DeleteChatByID, restricted)
	chatGroup.POST("/:chatID/messages", hc.CreateChatMessage, restricted)
}
