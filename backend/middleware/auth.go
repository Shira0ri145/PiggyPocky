package middleware

import (
	"os"
	"piggybackend/db"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func Protect() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing token"})
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}

		// ตรวจว่า token นี้ถูก revoke ใน redis หรือยัง
		exists, err := db.Redis.Get(db.Ctx, "blacklist:"+tokenString).Result()
		if err == nil && exists == "true" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "token revoked"})
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Locals("user_id", claims["user_id"])
		c.Locals("role", claims["role"])

		return c.Next()
	}
}

func OnlyAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "admin access only"})
		}
		return c.Next()
	}
}
