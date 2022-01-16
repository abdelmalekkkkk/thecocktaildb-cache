package api

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Config ...
type Config struct {
	Filter       func(*fiber.Ctx) bool
	ItemsPerPage int
}

// New ...
func New(config ...Config) fiber.Handler {
	var cfg Config
	if len(config) > 0 {
		cfg = config[0]
	}
	if cfg.ItemsPerPage == 0 {
		cfg.ItemsPerPage = 50
	}
	return func(c *fiber.Ctx) error {
		if cfg.Filter != nil && cfg.Filter(c) {
			return c.Next()
		}
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page <= 0 {
			page = 1
		}
		offset := cfg.ItemsPerPage * page
		c.Locals("offset", offset)
		c.Locals("limit", cfg.ItemsPerPage)
		return c.Next()
	}
}
