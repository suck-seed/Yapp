package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/suck-seed/yapp/internal/database"

	"github.com/gin-gonic/gin"
)

var (
	httpPort         string
	postgresUser     string
	postgresPass     string
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
		CORS:       corsMiddleware(),
		Postgres:   pgPool,
	}, nil
}

// corsMiddleware : Inject CORS settings into app
func corsMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

// loadEnvVariables : Loads env variables and injects them into var defined
func loadEnvVariables() error {

	err := godotenv.Load()
	if err != nil {
		fmt.Print("No .env File, proceeding with real ENV vars")
	}

	httpPort = os.Getenv("HTTP_PORT")
	postgresUser = os.Getenv("POSTGRES_USER")
	postgresPass = os.Getenv("POSTGRES_PASS")

	// since everything is inside docker container, it used default postgres port
	// db_host_port is for local app like TablePlus or pgAdmin to connect to postgres

	postgresHostPort = os.Getenv("HOST_POSTGRES_PORT")
	postgresHost = os.Getenv("POSTGRES_HOST")
	postgresPort = os.Getenv("POSTGRES_PORT")
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

	// TODO: Handle further error checking here

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
