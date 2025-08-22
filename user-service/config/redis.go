package config

import (
	"context"
	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

func NewRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis server address
	})

	// Test the connection
	if err := client.Ping(Ctx).Err(); err != nil {
		panic(err)
	}

	return client
}
