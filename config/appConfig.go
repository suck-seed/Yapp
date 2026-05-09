package config

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/suck-seed/yapp/internal/database"
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

// AppConfig : Stores configurations for server, includes port and db.
// CORS is handled by Nginx (see infra/nginx/nginx.conf), not by Gin.
type AppConfig struct {
	ServerPort   string
	CORS         gin.HandlerFunc
	PostgresPool *pgxpool.Pool
	RedisClient  *redis.Client
	RabbitMQURL  string
	NodeID       string
}

// SetupEnvironment : Loads ENV variables and returns the configurations
func SetupEnvironment() (config AppConfig, err error) {

	// loading environment variables
	err = loadEnvVariables()
	if err != nil {
		return AppConfig{}, err
	}

	// setting up gin mode
	if envMode := os.Getenv("APP_ENV"); envMode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	if endMode := os.Getenv("APP_ENV"); endMode == "development" {
		gin.SetMode(gin.DebugMode)
	}

	// load postgres instance
	pgPool, err := database.PostgresDBConnection()
	if err != nil {
		return AppConfig{}, err
	}

	// loading redis instance
	rdb, err := database.RedisDBConnection()
	if err != nil {
		return AppConfig{}, err
	}

	return AppConfig{
		ServerPort:   os.Getenv("PORT"),
		CORS:         buildCORS(),
		PostgresPool: pgPool,
		RedisClient:  rdb,
	}, nil
}
