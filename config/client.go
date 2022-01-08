package config

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

/*
	Connects to the Redis Server and returns a client
*/
func RedisClient(ctx *context.Context) (*redis.Client, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("There was a problem loading .env file")
	}

	redisHost := os.Getenv("REDIS_SERVER")
	redisPort := os.Getenv("REDIS_PORT")

	rdb := redis.NewClient(&redis.Options{
		Addr: redisHost + ":" + redisPort,
		Password: "",	
		DB:	0,
	})
	
	_, err = rdb.Ping(*ctx).Result()

	return rdb, err
}