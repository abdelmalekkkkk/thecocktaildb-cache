package api

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

// Controller has methods for fetching ingredients and cocktails from Redis
type Controller struct {
	Ctx   *context.Context
	Redis *redis.Client
}

// GetRecords fetches all records (ingredients, alcohols or cocktails) from Redis
func (controller Controller) GetRecords(c *fiber.Ctx) error {
	path := c.Path()
	var index string

	switch path {
	case "/ingredients":
		index = "idx:ingredients"
	case "/alcohols":
		index = "idx:alcohols"
	}

	results, err := controller.Redis.Do(*controller.Ctx, "FT.SEARCH", index, "*", "LIMIT", c.Locals("offset"), c.Locals("limit")).Slice()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"ingredients": formatRedisOutput(results)})
}

func formatRedisOutput(output []interface{}) []interface{} {

	var results []interface{}

	length := len(output)

	for i := 2; i < length; i += 2 {
		current := output[i].([]interface{})
		size := len(current)
		result := map[string]string{}
		for j := 0; j < size; j += 2 {
			key := current[j].(string)
			value := current[j+1].(string)
			result[key] = value
		}
		results = append(results, result)
	}

	return results

}
