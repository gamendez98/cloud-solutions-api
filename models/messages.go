package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Sender string

const (
	User      Sender = "user"
	Assistant Sender = "assistant"
	System    Sender = "system"
)

func (s *Sender) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	switch Sender(value) {
	case User, Assistant, System:
		*s = Sender(value)
		return nil
	default:
		return fmt.Errorf("invalid sender value: %s", value)
	}
}

type Message struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Sender    Sender    `json:"sender"`
	Text      string    `json:"text"`
}

func (chat *Chat) GetMessages() []Message {
	var messages []Message
	if !chat.Messages.Valid {
		panic(errors.New("chat does not have messages"))
	}
	err := json.Unmarshal(chat.Messages.RawMessage, &messages)
	if err != nil {
		panic(err)
	}
	return messages
}
