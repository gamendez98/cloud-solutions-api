package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	BucketName            string
	GCPServiceAccountFile string
	DbHost                string
	DbPort                string
	DbName                string
	DbUser                string
	DbPassword            string
	RabbitMQHost          string
	RabbitMQUsername      string
	RabbitMQPassword      string
	RabbitMQPort          string
	Host                  string
	Secret                string
	ProtocolPrefix        string
}

var config *Config

func GetConfig() *Config {
	if config != nil {
		return config
	}
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	config = new(Config)

	config.BucketName = os.Getenv("BUCKET_NAME")
	config.GCPServiceAccountFile = os.Getenv("GCP_SERVICE_ACCOUNT_FILE")
	config.DbHost = os.Getenv("DB_HOST")
	config.DbPort = os.Getenv("DB_PORT")
	config.DbName = os.Getenv("DB_NAME")
	config.DbPassword = os.Getenv("DB_PASSWORD")
	config.DbUser = os.Getenv("DB_USER")
	config.RabbitMQHost = os.Getenv("RABBIT_MQ_HOST")
	config.RabbitMQUsername = os.Getenv("RABBIT_MQ_USERNAME")
	config.RabbitMQPassword = os.Getenv("RABBIT_MQ_PASSWORD")
	config.RabbitMQPort = os.Getenv("RABBIT_MQ_PORT")
	config.Secret = os.Getenv("SECRET")
	config.ProtocolPrefix = "http"
	if os.Getenv("HTTPS") == "true" {
		config.ProtocolPrefix = "https"
	}

	fmt.Println(config)

	return config
}
