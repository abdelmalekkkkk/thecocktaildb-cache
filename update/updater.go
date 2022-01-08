package update

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/go-resty/resty/v2"
)

type Updater struct {
	Ctx *context.Context
	Redis *redis.Client
	Http *resty.Client
}



func (Updater) RunUpdate() {
	fmt.Println("Running update")
}