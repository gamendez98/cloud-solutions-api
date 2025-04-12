package handlers

import (
	"cloud-solutions-api/config"
	"cloud-solutions-api/models"
	"cloud-solutions-api/rabbitMQPublishers"
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/option"
	"strconv"
)

type HandlerContext struct {
	Queryer                     *models.Queries
	DocumentIndexingPublisher   *rabbitMQPublishers.DocumentIndexingPublisher
	AIAssistantMessagePublisher *rabbitMQPublishers.AIAssistantMessagePublisher
	StorageClient               *storage.Client
	Bucket                      *storage.BucketHandle
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
		fmt.Println(err)
		panic(err)
	}
	handlerContext.Queryer = queryer

	documentIndexingPublisher, err := rabbitMQPublishers.NewDocumentIndexingPublisher()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	handlerContext.DocumentIndexingPublisher = documentIndexingPublisher

	aiAssistantMessagePublisher, err := rabbitMQPublishers.NewAIAssistantMessagePublisher()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	handlerContext.AIAssistantMessagePublisher = aiAssistantMessagePublisher
	handlerContext.StorageClient, err = storage.NewClient(context.Background(), option.WithCredentialsFile(configuration.GCPServiceAccountFile))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	handlerContext.Bucket = handlerContext.StorageClient.Bucket(configuration.BucketName)

	return handlerContext
}

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
