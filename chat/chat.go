package chat

import (
	"cloud-solutions-api/models"
	"context"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"time"
)

type NewMessageParameters struct {
	Sender models.Sender `json:"sender"`
	Text   string        `json:"text"`
}

func NewMessageForChat(chat models.Chat, newMessage NewMessageParameters) models.Message {
	messages := chat.GetMessages()
	messageCount := len(messages)
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

func toGenAIMessage(message models.Message) *genai.Content {
	role := "user"
	if message.Sender == models.Assistant {
		role = "model"
	}
	return &genai.Content{
		Parts: []genai.Part{
			genai.Text(message.Text),
		},
		Role: role,
	}
}

func GetAssistantMessage(chat models.Chat, geminiModel *genai.GenerativeModel) (models.Message, error) {
	ctx := context.Background()

	messages := chat.GetMessages()
	chatSessions := geminiModel.StartChat()
	aiMessages := make([]*genai.Content, 0)
	for _, message := range messages[:len(messages)-1] {
		if len(message.Text) > 0 {
			aiMessages = append(aiMessages, toGenAIMessage(message))
		}
	}
	chatSessions.History = aiMessages

	resp, err := chatSessions.SendMessage(ctx, genai.Text(messages[len(messages)-1].Text))
	if err != nil {
		return models.Message{}, err
	}

	text := ""
	if resp.Candidates != nil {
		for _, v := range resp.Candidates[:1] {
			for _, k := range v.Content.Parts {
				partText := k.(genai.Text)
				text += string(partText)
			}
		}
	}

	return NewMessageForChat(
		chat,
		NewMessageParameters{
			Sender: models.Assistant,
			Text:   text,
		},
	), nil
}
