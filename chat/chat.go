package chat

import (
	"cloud-solutions-api/models"
	"fmt"
	"time"
)

type NewMessageParameters struct {
	Sender models.Sender `json:"sender"`
	Text   string        `json:"text"`
}

func NewMessageForChat(chat models.Chat, newMessage NewMessageParameters) models.Message {
	messageCount := len(chat.GetMessages())
	id := fmt.Sprintf("msg%d", messageCount)
	currentTime := time.Now()
	message := models.Message{
		ID:        id,
		Timestamp: currentTime,
		Sender:    newMessage.Sender,
		Text:      newMessage.Text,
	}
	return message
}
