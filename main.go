package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Loukay/thecokctaildb-cache/api"
	"github.com/Loukay/thecokctaildb-cache/config"
	"github.com/Loukay/thecokctaildb-cache/update"
	"github.com/go-co-op/gocron"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	var app *fiber.App = fiber.New(fiber.Config{
		Prefork: false,
	})

	controller := api.Controller{
		Redis: redis,
		Ctx:   &ctx,
	}

	app.Use(cors.New())
	app.Use(api.New())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON("The CocktailDB Cache")
	})

	app.Get("/ingredients", controller.GetRecords)
	app.Get("/alcohols", controller.GetRecords)

	err = app.Listen(":3001")

	if err != nil {
		log.Fatal("Failed to listen to web server")
		panic(err)
	}

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
