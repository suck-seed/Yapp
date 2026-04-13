package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// loadEnvVariables loads env variables and validates them based on APP_ENV.
func loadEnvVariables() error {
	err := godotenv.Load()
	if err != nil {
		fmt.Print("No .env file, proceeding with real ENV vars\n")
	}

	envMode := os.Getenv("APP_ENV")
	httpPort = os.Getenv("PORT")
	secretKey = os.Getenv("JWT_SECRET_KEY")

	if envMode != "production" && envMode != "development" {
		return errors.New("invalid environment mode")
	}

	if httpPort == "" {
		return errors.New("forgot to set PORT")
	}

	intPort, err := strconv.Atoi(httpPort)
	if err != nil {
		return errors.New("PORT must be a valid number")
	}
	if intPort <= 0 {
		return errors.New("PORT cannot be <= 0")
	}

	if secretKey == "" {
		return errors.New("forgot to set JWT_SECRET_KEY")
	}

	// Production validation
	if envMode == "production" {
		databaseURL := os.Getenv("DATABASE_URL")
		if databaseURL == "" {
			return errors.New("forgot to set DATABASE_URL")
		}
		return nil
	}

	// Development validation
	postgresHostPort = os.Getenv("HOST_POSTGRES_PORT")
	postgresHost = os.Getenv("POSTGRES_HOST")
	postgresPort = os.Getenv("POSTGRES_PORT")
	postgresUser = os.Getenv("POSTGRES_USER")
	postgresPassword = os.Getenv("POSTGRES_PASSWORD")
	postgresDbName = os.Getenv("POSTGRES_DB")

	if postgresPort == "" {
		postgresPort = "5432"
	}

	intPostgresPort, err := strconv.Atoi(postgresPort)
	if err != nil {
		return errors.New("POSTGRES_PORT must be a valid number")
	}
	if intPostgresPort <= 0 {
		return errors.New("POSTGRES_PORT cannot be <= 0")
	}

	if postgresHost == "" {
		return errors.New("forgot to set POSTGRES_HOST")
	}

	if postgresUser == "" || postgresPassword == "" || postgresDbName == "" {
		return errors.New("missing required DB environment variables: POSTGRES_USER/POSTGRES_PASSWORD/POSTGRES_DB")
	}

	return nil
}
