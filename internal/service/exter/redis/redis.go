package redisservice

import (
	"context"
	"fmt"
	"log"

	redis "github.com/go-redis/redis/v8"
)

var RDB *redis.Client

func New() error {
	// Create a context for the Redis client operations
	ctx := context.Background()

	// Connect to Redis (default address is localhost:6379)
	RDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // No password set
		DB:       0,                // Default DB
	})

	// Test the connection
	err := RDB.Ping(ctx).Err()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	fmt.Println("Connected to Redis successfully!")

	return nil
}
