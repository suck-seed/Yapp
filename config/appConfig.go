package config

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"os"
)

type AppConfig struct {
	ServerPort string
	CORS       gin.HandlerFunc
}

// Main function
func SetupEnvironment() (config AppConfig, err error) {

	// load the env file

	if os.Getenv("APP_ENV") == "dev" {

		// load the local .env file
		err = godotenv.Load()

		// if file doesnt exist
		if err != nil {
			return AppConfig{}, err

			//! FOR PRODUCTION, REMOVE RETURN
			//log.Println("⚠️  no .env file found, proceeding with real ENV vars")

		}

	}

	// get ENV variables
	httpPort := os.Getenv("HTTP_PORT")

	if len(httpPort) < 1 {
		return AppConfig{}, errors.New("Forgot to set HTTP_PORT ? ")
	}

	return AppConfig{
		ServerPort: httpPort,
		CORS:       corsMiddleware(),
	}, nil
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
