package config

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
)

type AppConfig struct {
	ServerPort string
}

// Main func
func SetupEnvironment() (config AppConfig, err error) {

	if os.Getenv("APP_ENV") != "production" {
		godotenv.Load()
	}

	//TODO populate env file and implement this
	httpPort := os.Getenv("HTTP_PORT")

	if httpPort == "" || len(httpPort) < 1 {
		return AppConfig{}, errors.New("Forgot to set HTTP_PORT ? ")
	}

	return AppConfig{
		ServerPort: httpPort,
	}, nil
}
