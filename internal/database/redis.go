package database

import (
	"errors"
	"os"

	"github.com/redis/go-redis/v9"
)

func RedisDBConnection() (*redis.Client, error) {

	// Read environment variables directly
	redisAddress := os.Getenv("REDIS_ADDRESS")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	if redisAddress == "" || redisPassword == "" {

		return nil, errors.New("Empty redis env fields")
	}

	return redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: redisPassword,
	}), nil
}
