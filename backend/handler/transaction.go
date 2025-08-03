package handlers

import (
	"coin-backend/models"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func GetTransactions(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	return c.JSON(models.GetAllTransactions(userID))
}

func GetTransactionByID(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	tx, err := models.GetTransactionByID(id)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	return c.JSON(tx)
}
