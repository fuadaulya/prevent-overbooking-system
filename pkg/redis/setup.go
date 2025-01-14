package redis

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

func InitRedis() *redis.Client {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	if redisHost == "" {
		redisHost = "localhost" // Default to localhost if not found in environment
	}

	redisAddr := redisHost + ":" + redisPort

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr, // Use the configured address
		Password: "",        // No password by default
		DB:       0,         // Use default DB
	})

	// Test Redis connection
	err := client.Ping(context.Background()).Err()
	if err != nil {
		log.Fatalf("failed to connect to Redis: %+v", err)
	}
	log.Println("Connected to Redis successfully")
	//_________________

	return client
}
