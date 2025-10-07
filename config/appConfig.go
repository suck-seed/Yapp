package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
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

// corsMiddleware : Inject CORS settings into app
func buildCORS() gin.HandlerFunc {

	_ = godotenv.Load()

	origin := os.Getenv("FRONTEND_ORIGIN")
	if origin == "" {
		origin = "http://localhost:3000"
	}

	cfg := cors.Config{
		AllowOrigins:     []string{origin},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "X-CSRF-Token", "Origin"},
		ExposeHeaders:    []string{"Set-Cookie"},
		AllowCredentials: true, // <- required for cookies
		MaxAge:           12 * time.Hour,
	}
	cfg.AddAllowHeaders("Vary")
	return cors.New(cfg)
}

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

func GetSecretKey() string {

	err := godotenv.Load()
	if err != nil {
		fmt.Print("No .env File, proceeding with real ENV vars")
	}

	secretKey = string(os.Getenv("JWT_SECRET_KEY"))
	return secretKey

}

func FrontEndOrigin() string {

	err := godotenv.Load()
	if err != nil {
		fmt.Print("No .env File, proceeding with real ENV vars")
	}

	origin := string(os.Getenv("FRONTEND_ORIGIN"))
	return origin

}
