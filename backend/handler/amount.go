package handlers

import (
	"coin-backend/models"
	"github.com/gofiber/fiber/v2"
)

func Amount(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	amount, err := models.GetAmountByUserID(userID)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	return c.JSON(amount)
}
