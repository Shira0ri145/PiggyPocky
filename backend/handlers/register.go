package handlers

import (
	"piggybackend/db"
	"piggybackend/models"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// Register godoc
// @Summary Register a new user (Admin only)
// @Description Create a new user (admin only). Also creates amount with default value = 0
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param register body models.RegisterInput true "New user info"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /api/register [post]
func Register(c *fiber.Ctx) error {
	var input models.RegisterInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	user := models.User{
		Username:     input.Username,
		PasswordHash: string(hash),
		Role:         input.Role,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot create user"})
	}

	// ใช้ Value แทน Balance
	amount := models.Amount{UserID: user.ID, Balance: 0.00}
	db.DB.Create(&amount)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "user created"})
}
