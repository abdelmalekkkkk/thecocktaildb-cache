package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
)

func main() {

	s := gocron.NewScheduler(time.UTC)

	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("There was a problem loading .env file")
	}

	redis, err := RedisClient(&ctx)

	if err != nil {
		panic("There was a problem connecting to the Redis server.")
	}

	updater := Updater{
		Ctx:   &ctx,
		Redis: redis,
		API:   NewAPIClient(os.Getenv("API_URL")),
	}

	s.Every(6).Hours().Do(updater.Run)

	s.StartBlocking()

}
