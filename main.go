package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Loukay/thecokctaildb-cache/config"
	"github.com/Loukay/thecokctaildb-cache/update"
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

	redis, err := config.RedisClient(&ctx)

	if err != nil {
		panic("There was a problem connecting to the Redis server.")
	}

	updater := update.Updater{
		Ctx:   &ctx,
		Redis: redis,
		API:   update.NewAPIClient(os.Getenv("API_URL")),
	}

	s.Every(6).Hours().Do(updater.Run)

	s.StartBlocking()

}
