package config

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	httpPort string
	dbUser   string
	dbPass   string
	dbHost   string
	dbPort   string
)

// AppConfig : Stores configurations for server, includes port, db, and middleware
type AppConfig struct {
	ServerPort string
	CORS       gin.HandlerFunc
}

// SetupEnvironment : Loads ENV variables, creates instance of middleware and returns the configurations
func SetupEnvironment() (config AppConfig, err error) {

	// Load Env Variables
	err = loadEnvVariables()
	if err != nil {
		return AppConfig{}, err
	}

	return AppConfig{
		ServerPort: httpPort,
		CORS:       corsMiddleware(),
	}, nil
}

// loadEnvVariables : Loads env variables and injects them into var defined
func loadEnvVariables() error {

	err := godotenv.Load()
	if err != nil {
		fmt.Print("No .env File, proceeding with real ENV vars")
	}

	// get ENV variables
	httpPort = os.Getenv("HTTP_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort = os.Getenv("DB_PORT")

	if len(httpPort) < 1 {
		return errors.New("Forgot to set HTTP_PORT ? ")
	}

	if len(dbPort) < 1 {
		return errors.New("Forgot to set DB_PORT ? ")
	}

	return nil

}

// CORS middleware
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

func createDBConnection() (error, *sql.DB) {

	return nil, &sql.DB{}
}
