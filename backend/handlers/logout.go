package handlers

import (
	"time"

	"piggybackend/db"

	"github.com/gofiber/fiber/v2"
)

// Logout godoc
// @Summary Logout current user
// @Description Revoke JWT token by adding it to Redis blacklist
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /logout [post]
func Logout(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing token"})
	}
	tokenString = tokenString[len("Bearer "):]

	// บันทึก token ที่ revoke ไว้ใน redis
	expiration := time.Hour * 24 // TTL เท่ากับอายุ token
	err := db.Redis.Set(db.Ctx, "blacklist:"+tokenString, "true", expiration).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not revoke token"})
	}

	return c.JSON(fiber.Map{"message": "logout successful"})
}
