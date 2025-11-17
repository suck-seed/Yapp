package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func GetSecretKey() string {

	err := godotenv.Load()
	if err != nil {
		fmt.Print("No .env File, proceeding with real ENV vars")
	}

	secretKey = string(os.Getenv("JWT_SECRET_KEY"))
	return secretKey

}
