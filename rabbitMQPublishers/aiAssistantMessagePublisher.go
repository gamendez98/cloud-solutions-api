package rabbitMQPublishers

import "cloud-solutions-api/models"

type AIAssistantMessagePublisher struct {
	RabbitMQPublisher
}

type AIAssistantMessage struct {
	Messages []models.Message `json:"messages"`
}

const AIAssistantExchangeName = "ai-assistant-exchange"
const AIAssistantQueueName = "ai-assistant-queue"

func NewAIAssistantMessagePublisher() (*AIAssistantMessagePublisher, error) {
	rabbitMQPublisher, err := NewRabbitMQPublisher(
		AIAssistantExchangeName,
		AIAssistantQueueName,
	)
	if err != nil {
		return nil, err
	}
	return &AIAssistantMessagePublisher{
		RabbitMQPublisher: *rabbitMQPublisher,
	}, nil
}

func (aip *AIAssistantMessagePublisher) Publish(message AIAssistantMessage) error {
	return aip.PublishJSON(message)
}
