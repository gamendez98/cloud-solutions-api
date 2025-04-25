package pubSubPublisher

import (
	"cloud-solutions-api/models"
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"
	"google.golang.org/api/option"
)

type PubSubPublisher struct {
	Client                *pubsub.Client
	DocumentIndexingTopic *pubsub.Topic
	AiAssistantTopic      *pubsub.Topic
}

type DocumentIndexingMessage struct {
	DocumentId   int32  `json:"document_id"`
	DocumentText string `json:"document_text"`
}

type AIAssistantMessage struct {
	ChatId   int32            `json:"chat_id"`
	Messages []models.Message `json:"messages"`
}

const DocumentIndexingTopicName = "DocumentIndexing"
const AiAssistantTopicName = "AiAssistant"

func NewPubSubPublisher(gcpProjectID string, gcpServiceAccountFile string) (*PubSubPublisher, error) {
	client, err := pubsub.NewClient(context.Background(), gcpProjectID, option.WithCredentialsFile(gcpServiceAccountFile))
	if err != nil {
		return nil, err
	}
	documentIndexingTopic := client.Topic(DocumentIndexingTopicName)
	aiAssistantTopic := client.Topic(AiAssistantTopicName)
	publisher := &PubSubPublisher{
		Client:                client,
		DocumentIndexingTopic: documentIndexingTopic,
		AiAssistantTopic:      aiAssistantTopic,
	}
	return publisher, nil
}

func (publisher *PubSubPublisher) Close() error {
	publisher.DocumentIndexingTopic.Stop()
	publisher.AiAssistantTopic.Stop()
	if err := publisher.Client.Close(); err != nil {
		return err
	}
	return nil
}

func publishJson(topic *pubsub.Topic, message interface{}) error {
	ctx := context.Background()
	marshalled, err := json.Marshal(message)
	if err != nil {
		return err
	}
	result := topic.Publish(ctx, &pubsub.Message{
		Data: marshalled,
		Attributes: map[string]string{
			"origin": "api",
		},
	})

	go func() {
		_, err = result.Get(ctx)
		if err != nil {
			log.Error(fmt.Sprintf("Error publishing message to %s: %s", topic.String(), err))
		}
	}()

	return nil
}

func (publisher *PubSubPublisher) PublishDocumentIndexingMessage(message DocumentIndexingMessage) error {
	return publishJson(publisher.DocumentIndexingTopic, message)
}

func (publisher *PubSubPublisher) PublishAiAssistantMessage(message AIAssistantMessage) error {
	return publishJson(publisher.AiAssistantTopic, message)
}
