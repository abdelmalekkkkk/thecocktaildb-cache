package main

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

// RedisClient connects to the Redis Server and returns a client
func RedisClient(ctx *context.Context) (*redis.Client, error) {
	redisHost := os.Getenv("REDIS_SERVER")
	redisPort := os.Getenv("REDIS_PORT")

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping(*ctx).Result()

	return rdb, err
}
