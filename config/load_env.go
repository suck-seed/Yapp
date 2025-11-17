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

	httpPort = os.Getenv("HTTP_PORT")
	postgresUser = os.Getenv("POSTGRES_USER")
	postgresPassword = os.Getenv("POSTGRES_PASSWORD")

	// since everything is inside docker container, it used default postgres port
	// db_host_port is for local app like TablePlus or pgAdmin to connect to postgres

	postgresHostPort = os.Getenv("HOST_POSTGRES_PORT")
	postgresHost = os.Getenv("POSTGRES_HOST")
	if postgresHost == "" {
		// default for docker compose network; outside docker use localhost
		if os.Getenv("APP_ENV") == "dev" {
			postgresHost = "postgres"
		} else {
			postgresHost = "localhost"
		}
	}
	postgresPort = os.Getenv("POSTGRES_PORT")
	if postgresPort == "" {
		postgresPort = "5432"
	}
	postgresDbName = os.Getenv("POSTGRES_DB")
	secretKey = os.Getenv("JWT_SECRET_KEY")

	// HANDLING INAPPROPRIATE ENV VARIABLES

	intPort, err := strconv.Atoi(httpPort)
	if len(httpPort) < 1 {

		if intPort <= 0 {
			return errors.New("Http Port cannot be <= 0s")
		}

		return errors.New("Forgot to set HTTP_PORT ? ")
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
