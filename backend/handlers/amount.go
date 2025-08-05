package handlers

import (
	"piggybackend/db"
	"piggybackend/models"

	"github.com/gofiber/fiber/v2"
)

// Amount godoc
// @Summary Get amount of current user
// @Description Get current user's balance amount
// @Tags Amount
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.AmountResponse
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/amount [get]
func Amount(c *fiber.Ctx) error {
	// ดึงค่า user_id จาก JWT (Locals เก็บเป็น float64 เสมอ)
	userIDFloat, ok := c.Locals("user_id").(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user_id"})
	}
	userID := uint(userIDFloat)

	amount, err := models.GetAmountByUserID(db.DB, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "amount not found"})
	}
	return c.JSON(models.AmountResponse{
		UserID:  userID,
		Balance: amount.Balance,
	})
}
