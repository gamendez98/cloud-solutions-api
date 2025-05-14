package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	BucketName            string
	GCPServiceAccountFile string
	GCPProjectID          string
	DbHost                string
	DbPort                string
	DbName                string
	DbUser                string
	DbPassword            string
	Host                  string
	Secret                string
	ProtocolPrefix        string
	Port                  string
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
	config.GCPServiceAccountFile = "api-service-account.json"
	config.GCPProjectID = os.Getenv("GCP_PROJECT_ID")
	config.DbHost = os.Getenv("DB_HOST")
	config.DbPort = os.Getenv("DB_PORT")
	config.DbName = os.Getenv("DB_NAME")
	config.DbPassword = os.Getenv("DB_PASSWORD")
	config.DbUser = os.Getenv("DB_USER")
	config.Secret = os.Getenv("SECRET")
	config.Port = os.Getenv("PORT")
	if config.Port == "" {
		config.Port = "80"
	}

	fmt.Println(config)

	return config
}
