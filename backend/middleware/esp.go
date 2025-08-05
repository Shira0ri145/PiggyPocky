package middleware

import (
	"piggybackend/config"

	"github.com/gofiber/fiber/v2"
)

func OnlyESP() fiber.Handler {
	return func(c *fiber.Ctx) error {
		key := c.Get("X-ESP-KEY")
		if key == "" || key != config.GetESPApiKey() {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized device"})
		}
		return c.Next()
	}
}
