package handlers

import (
	"cloud-solutions-api/config"
	"cloud-solutions-api/models"
	"cloud-solutions-api/pubSubPublisher"
	"cloud.google.com/go/storage"
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"google.golang.org/api/option"
	"strconv"
)

type HandlerContext struct {
	Queryer        *models.Queries
	PuSubPublisher *pubSubPublisher.PubSubPublisher
	StorageClient  *storage.Client
	Bucket         *storage.BucketHandle
	Secret         []byte
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
		log.Error(err)
		panic(err)
	}
	handlerContext.Queryer = queryer
	handlerContext.StorageClient, err = storage.NewClient(context.Background(), option.WithCredentialsFile(configuration.GCPServiceAccountFile))
	handlerContext.Bucket = handlerContext.StorageClient.Bucket(configuration.BucketName)
	publisher, err := pubSubPublisher.NewPubSubPublisher(configuration.GCPProjectID, configuration.GCPServiceAccountFile)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	handlerContext.PuSubPublisher = publisher

	return handlerContext
}

func (hc *HandlerContext) HealthCheck(c echo.Context) error {
	return c.String(200, "OK")
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

func (hc *HandlerContext) Close() []error {
	var errs []error
	err := hc.PuSubPublisher.Close()
	errs = append(errs, err)
	err = hc.StorageClient.Close()
	errs = append(errs, err)
	return errs
}
