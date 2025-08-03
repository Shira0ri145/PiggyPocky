package handlers

import (
	"coin-backend/models"
	"github.com/gofiber/fiber/v2"
)

func Profile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	user, err := models.GetUserByID(userID)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	return c.JSON(user)
}
