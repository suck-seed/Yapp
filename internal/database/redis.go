package database

import (
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RedisDBConnection() (*redis.Client, error) {

	// FOR PRODUCTION
	if gin.Mode() == gin.ReleaseMode {

		redisURL := os.Getenv("REDIS_URL")
		if redisURL == "" {
			panic("REDIS_URL is not set")
		}

		// parsing url into opts
		redisOpts, err := redis.ParseURL(redisURL)
		if err != nil {
			return nil, fmt.Errorf("REDIS_URL not set")
		}

		// creating new redis client
		rdb := redis.NewClient(redisOpts)

		// test connection, ping
		_, err = rdb.Ping(context.Background()).Result()
		if err != nil {
			return nil, fmt.Errorf("Failed to establish an conection to Redis")
		}

		return rdb, nil
	} else if gin.Mode() == gin.DebugMode {

		redisURL := os.Getenv("REDIS_URL")
		if redisURL == "" {
			panic("REDIS_URL is not set")
		}

		// parsing url into opts
		redisOpts, err := redis.ParseURL(redisURL)
		if err != nil {
			return nil, fmt.Errorf("REDIS_URL not set")
		}

		// creating new redis client
		rdb := redis.NewClient(redisOpts)

		// test connection, ping
		_, err = rdb.Ping(context.Background()).Result()
		if err != nil {
			return nil, fmt.Errorf("Failed to establish an conection to Redis")
		}

		return rdb, nil

	}

	return nil, fmt.Errorf("gin mode is not set")
}
