package handlers

import (
	"cloud-solutions-api/config"
	"cloud-solutions-api/models"
	"cloud-solutions-api/rabbitConnection"
)

type HandlerContext struct {
	Queryer                     *models.Queries
	DocumentIndexingPublisher   *rabbitConnection.DocumentIndexingPublisher
	AIAssistantMessagePublisher *rabbitConnection.AIAssistantMessagePublisher
	Secret                      []byte
}

func NewHandlerContext(configuration config.Config) *HandlerContext {
	handlerContext := &HandlerContext{
		Secret: []byte(configuration.Secret),
	}
	queryer, err := models.NewQueryer(models.Config{
		DBHost:     configuration.DbHost,
		DBPort:     configuration.DbPort,
		DBUser:     configuration.DbUser,
		DBName:     configuration.DbName,
		DBPassword: configuration.DbPassword,
	})
	if err != nil {
		panic(err)
	}
	handlerContext.Queryer = queryer

	documentIndexingPublisher, err := rabbitConnection.NewDocumentIndexingPublisher()
	if err != nil {
		panic(err)
	}
	handlerContext.DocumentIndexingPublisher = documentIndexingPublisher

	aiAssistantMessagePublisher, err := rabbitConnection.NewAIAssistantMessagePublisher()
	if err != nil {
		panic(err)
	}
	handlerContext.AIAssistantMessagePublisher = aiAssistantMessagePublisher

	return handlerContext
}
