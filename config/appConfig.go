package config

import (
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suck-seed/yapp/internal/database"

	"github.com/gin-gonic/gin"
)

var (
	httpPort         string
	postgresUser     string
	postgresPassword string
	postgresHost     string
	postgresHostPort string
	postgresPort     string
	postgresDbName   string
	secretKey        string
)

// AppConfig : Stores configurations for server, includes port, db, and middleware
type AppConfig struct {
	ServerPort string
	CORS       gin.HandlerFunc
	Postgres   *pgxpool.Pool
}

// SetupEnvironment : Loads ENV variables, creates instance of middleware and returns the configurations
func SetupEnvironment() (config AppConfig, err error) {

	// Load Env Variables
	err = loadEnvVariables()
	if err != nil {
		return AppConfig{}, err
	}

	// load postgres instance
	pgPool, err := database.PostgresDBConnection()
	if err != nil {
		return AppConfig{}, err
	}

	return AppConfig{
		ServerPort: os.Getenv("HTTP_PORT"),
		CORS:       buildCORS(),
		Postgres:   pgPool,
	}, nil
}
