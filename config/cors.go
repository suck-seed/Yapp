package config

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

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

	return cors.New(cfg)
}

func FrontEndOrigin() string {

	err := godotenv.Load()
	if err != nil {
		fmt.Print("No .env File, proceeding with real ENV vars")
	}

	origin := string(os.Getenv("FRONTEND_ORIGIN"))
	return origin

}
