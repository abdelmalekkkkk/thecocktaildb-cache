package main

import (
	"context"

	"github.com/Loukay/thecokctaildb-cache/config"
	"github.com/Loukay/thecokctaildb-cache/update"
)



func main() {

	ctx := context.Background()

	redis, err := config.RedisClient(&ctx)

	updater := update.Updater{
		Ctx: &ctx,
		Redis: redis,
		API: update.NewAPIClient(),
	}

	updater.Run()

	if err != nil {
		panic("There was a problem connecting to the Redis server.")
	}

}