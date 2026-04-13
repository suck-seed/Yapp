package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// loadEnvVariables : Loads env variables and injects them into var defined
func loadEnvVariables() error {

	err := godotenv.Load()
	if err != nil {
		fmt.Print("No .env File, proceeding with real ENV vars")
	}

	envMode := os.Getenv("APP_ENV")
	httpPort = os.Getenv("PORT")

	postgresHostPort = os.Getenv("HOST_POSTGRES_PORT")
	postgresHost = os.Getenv("POSTGRES_HOST")
	postgresPort = os.Getenv("POSTGRES_PORT")
	postgresUser = os.Getenv("POSTGRES_USER")
	postgresPassword = os.Getenv("POSTGRES_PASSWORD")
	postgresDbName = os.Getenv("POSTGRES_DB")

	secretKey = os.Getenv("JWT_SECRET_KEY")

	if postgresPort == "" {
		postgresPort = "5432"
	}

	// HANDLING INAPPROPRIATE ENV VARIABLES
	if envMode != "production" && envMode != "development" {
		return errors.New("Invalid environment mode")
	}

	intPort, err := strconv.Atoi(httpPort)
	if len(httpPort) < 1 {
		if intPort <= 0 {
			return errors.New("Http Port cannot be <= 0s")
		}
		return errors.New("Forgot to set PORT ? ")
	}

	intPostgresPort, err := strconv.Atoi(postgresPort)
	if len(postgresPort) < 1 {
		if intPostgresPort <= 0 {
			return errors.New("Postgres Port cannot be <= 0s")
		}
		return errors.New("Forgot to set POSTGRES_PORT ? ")
	}

	// Further validation: ensure critical DB vars exist
	if postgresUser == "" || postgresPassword == "" || postgresDbName == "" {
		return errors.New("Missing required DB environment variables: POSTGRES_USER/POSTGRES_PASSWORD/POSTGRES_DB")
	}
	return nil
}
