package handlers

import (
	"cloud-solutions-api/config"
	"cloud-solutions-api/models"
)

type HandlerContext struct {
	Queryer *models.Queries
	Secret  string
}

func NewHandlerContext(configuration config.Config) *HandlerContext {
	handlerContext := &HandlerContext{
		Secret: configuration.Secret,
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
	return handlerContext
}
