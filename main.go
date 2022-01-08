package main

import (
	"context"

	"github.com/Loukay/thecokctaildb-cache/config"
	"github.com/Loukay/thecokctaildb-cache/update"
	"github.com/go-resty/resty/v2"
)



func main() {

	ctx := context.Background()

	redis, err := config.RedisClient(&ctx)

	http := resty.New()

	updater := update.Updater{
		Ctx: &ctx,
		Redis: redis,
		Http: http,
	}

	updater.RunUpdate()

	if err != nil {
		panic("There was a problem connecting to the Redis server.")
	}

}