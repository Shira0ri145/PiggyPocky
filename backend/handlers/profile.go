package handlers

import (
	"piggybackend/db"
	"piggybackend/models"

	"github.com/gofiber/fiber/v2"
)

// Profile godoc
// @Summary Get current user profile
// @Description Get profile of the authenticated user
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.UserProfile
// @Failure 401 {object} map[string]string
// @Router /api/profile [get]
func Profile(c *fiber.Ctx) error {
	// ป้องกัน panic ด้วย type assertion ตรวจสอบก่อน
	userIDFloat, ok := c.Locals("user_id").(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid user_id in token",
		})
	}
	userID := uint(userIDFloat)

	user, err := models.GetUserByID(db.DB, userID)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	amount, _ := models.GetAmountByUserID(db.DB, userID)
	return c.JSON(models.UserProfile{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
		Balance:  amount.Balance,
	})
}
