package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"time"

	"github.com/Loukay/thecokctaildb-cache/config"
	"github.com/Loukay/thecokctaildb-cache/update"
	"github.com/go-co-op/gocron"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

func main() {

	s := gocron.NewScheduler(time.UTC)

	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("There was a problem loading .env file")
	}

	redis, err := config.RedisClient(&ctx)

	if err != nil {
		panic("There was a problem connecting to the Redis server.")
	}

	buildIndexes(&ctx, redis)

	updater := update.Updater{
		Ctx:   &ctx,
		Redis: redis,
		API:   update.NewAPIClient(os.Getenv("API_URL")),
	}

	s.Every(6).Hours().Do(updater.Run)

	s.StartAsync()

	bufio.NewScanner(os.Stdin).Scan()

}

func buildIndexes(ctx *context.Context, redis *redis.Client) {
	_, err :=
		redis.Do(*ctx, "FT.CREATE", "idx:ingredients",
			"ON", "hash",
			"PREFIX", "1", "ingredient:",
			"SCHEMA",
			"name", "TEXT",
			"type", "TEXT").Result()

	if err != nil {
		log.Printf("Couldn't create ingredients index %v", err)
	}

	_, err =
		redis.Do(*ctx, "FT.CREATE", "idx:alcohols",
			"ON", "hash",
			"PREFIX", "1", "alcohol:",
			"SCHEMA",
			"name", "TEXT",
			"type", "TEXT").Result()

	if err != nil {
		log.Printf("Couldn't create alcohols index %v", err)
	}

	_, err =
		redis.Do(*ctx, "FT.CREATE", "idx:cocktails",
			"ON", "hash",
			"PREFIX", "1", "cocktail:",
			"SCHEMA",
			"name", "TEXT",
			"category", "TEXT",
			"ingredients", "TAG").Result()
	if err != nil {
		log.Printf("Couldn't create Index %v", err)
	}
}
