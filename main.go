package main

import (
	"bufio"
	"context"
	"log"
	"os"

	"github.com/Loukay/thecokctaildb-cache/config"
	"github.com/Loukay/thecokctaildb-cache/update"
	"github.com/joho/godotenv"
)

func main() {

	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("There was a problem loading .env file")
	}

	redis, err := config.RedisClient(&ctx)

	if err != nil {
		panic("There was a problem connecting to the Redis server.")
	}

	res, err :=
		redis.Do(ctx, "FT.CREATE", "idx:ingredients",
			"ON", "hash",
			"PREFIX", "1", "ingredient:",
			"SCHEMA",
			"name", "TEXT",
			"type", "TEXT").Result()

	if err != nil {
		log.Printf("Couldn't create ingredients index %v", err)
	}

	res, err =
		redis.Do(ctx, "FT.CREATE", "idx:alcohols",
			"ON", "hash",
			"PREFIX", "1", "alcohol:",
			"SCHEMA",
			"name", "TEXT",
			"type", "TEXT").Result()

	if err != nil {
		log.Printf("Couldn't create alcohols index %v", err)
	}

	res, err =
		redis.Do(ctx, "FT.CREATE", "idx:cocktails",
			"ON", "hash",
			"PREFIX", "1", "cocktail:",
			"SCHEMA",
			"name", "TEXT",
			"category", "TEXT",
			"ingredients", "TAG").Result()
	if err != nil {
		log.Printf("Couldn't create Index %v", err)
	}

	log.Printf("Creating index: %v", res)

	updater := update.Updater{
		Ctx:   &ctx,
		Redis: redis,
		API:   update.NewAPIClient(os.Getenv("API_URL")),
	}

	updater.Run()

	bufio.NewScanner(os.Stdin).Scan()

}
