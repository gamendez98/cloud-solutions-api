package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	DbHost         string
	DbPort         string
	DbName         string
	DbUser         string
	DbPassword     string
	Host           string
	Secret         string
	ProtocolPrefix string
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

	config.DbHost = os.Getenv("DB_HOST")
	config.DbPort = os.Getenv("DB_PORT")
	config.DbName = os.Getenv("DB_NAME")
	config.DbPassword = os.Getenv("DB_PASSWORD")
	config.DbUser = os.Getenv("DB_USER")
	config.Secret = os.Getenv("SECRET")
	config.ProtocolPrefix = "http"
	if os.Getenv("HTTPS") == "true" {
		config.ProtocolPrefix = "https"
	}

	fmt.Println(config)

	return config
}
