package handlers

import (
	"fmt"
	"piggybackend/config"
	"piggybackend/db"
	"piggybackend/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// Login godoc
// @Summary Login
// @Description Authenticate user and return JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body models.LoginInput true "Login credentials"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /login [post]
func Login(c *fiber.Ctx) error {
	var input models.LoginInput
	secret := config.GetEnv("JWT_SECRET", "secret")

	// แปลง JSON body
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid input",
		})
	}

	// หาผู้ใช้จาก DB
	user, err := models.FindUserByUsername(db.DB, input.Username)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid credentials",
		})
	}
	fmt.Println("Hash in DB:", user.PasswordHash)
	fmt.Println("Input password:", input.Password)

	// ตรวจรหัสผ่าน bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid credentials",
		})
	}

	// สร้าง JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	t, err := token.SignedString([]byte(secret)) // 👉 ควรย้ายไป .env ในระบบจริง
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "could not generate token",
		})
	}

	// ส่ง JWT กลับ
	return c.JSON(fiber.Map{
		"token": t,
	})

}
