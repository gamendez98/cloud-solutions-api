package handlers

import (
	"cloud-solutions-api/authentication"
	"cloud-solutions-api/chat"
	"cloud-solutions-api/models"
	"cloud-solutions-api/pubSubPublisher"
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

// ChatOwnershipMiddleware checks if the current user owns the chat specified in the request
func (hc *HandlerContext) ChatOwnershipMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		chatIDString := c.Param("chatID")
		chatID, err := strconv.Atoi(chatIDString)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid chat ID")
		}

		if !CheckChatOwnership(hc, c, chatID) {
			return echo.NewHTTPError(http.StatusForbidden, "Forbidden: You do not own this chat")
		}

		// Continue to the next handler if ownership is verified
		return next(c)
	}
}

func (hc *HandlerContext) CreateEmptyChat(c echo.Context) error {
	account, err := authentication.GetCurrentAccount(hc.Queryer, c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
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
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid chat ID")
	}

	retrievedChat, err := hc.Queryer.GetChatByID(context.Background(), int32(chatID))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid chat ID")
	}

	return c.JSON(http.StatusOK, retrievedChat)
}

func (hc *HandlerContext) DeleteChatByID(c echo.Context) error {
	chatIDString := c.Param("chatID")
	chatID, err := strconv.Atoi(chatIDString)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid chat ID")
	}

	err = hc.Queryer.DeleteChat(context.Background(), int32(chatID))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid chat ID")
	}

	return c.JSON(http.StatusOK, echo.Map{})
}

func (hc *HandlerContext) CreateChatMessage(c echo.Context) error {
	chatIDString := c.Param("chatID")
	chatID, err := strconv.Atoi(chatIDString)
	newMessageParameters := chat.NewMessageParameters{}

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid chat ID")
	}

	if err = c.Bind(&newMessageParameters); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	retrievedChat, err := hc.Queryer.GetChatByID(context.Background(), int32(chatID))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid chat ID")
	}

	newMessage := chat.NewMessageForChat(
		retrievedChat,
		newMessageParameters,
	)
	newMessageJSON, err := json.Marshal(newMessage)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	retrievedChat, err = hc.Queryer.AddMessageToChat(context.Background(), models.AddMessageToChatParams{
		Chatid:     retrievedChat.ID,
		Newmessage: newMessageJSON,
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	err = hc.PuSubPublisher.PublishAiAssistantMessage(pubSubPublisher.AIAssistantMessage{
		Messages: retrievedChat.GetMessages(),
		ChatId:   retrievedChat.ID,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, retrievedChat)
}

func (hc *HandlerContext) GetChatIsUnread(c echo.Context) error {
	chatIDString := c.Param("chatID")
	chatID, err := strconv.Atoi(chatIDString)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid chat ID")
	}

	unread, err := hc.Queryer.IsUnread(context.Background(), int32(chatID))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid chat ID")
	}

	return c.JSON(http.StatusOK, map[string]bool{"unread": unread.Bool})
}

func (hc *HandlerContext) ChatMarkAsRead(c echo.Context) error {
	chatIDString := c.Param("chatID")
	chatID, err := strconv.Atoi(chatIDString)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid chat ID")
	}

	err = hc.Queryer.MarkAsReadByID(context.Background(), int32(chatID))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid chat ID")
	}

	return c.JSON(http.StatusOK, echo.Map{})
}

func RegisterChatRoutes(e *echo.Echo, hc *HandlerContext) {
	restricted := echojwt.JWT(hc.Secret)
	chatGroup := e.Group("/chats")
	chatGroup.GET("/:chatID", hc.GetChatByID, restricted, hc.ChatOwnershipMiddleware)
	chatGroup.GET("/:chatID/unread", hc.GetChatIsUnread, restricted, hc.ChatOwnershipMiddleware)
	chatGroup.POST("", hc.CreateEmptyChat, restricted)
	chatGroup.DELETE("/:chatID", hc.DeleteChatByID, restricted, hc.ChatOwnershipMiddleware)
	chatGroup.POST("/:chatID/messages", hc.CreateChatMessage, restricted, hc.ChatOwnershipMiddleware)
	chatGroup.POST("/:chatID/mark-as-read", hc.ChatMarkAsRead, restricted, hc.ChatOwnershipMiddleware)
}
