package middleware

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func Protect() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte("secret"),
		ContextKey:   "user",
		SuccessHandler: func(c *fiber.Ctx) error {
			claims := c.Locals("user").(*jwtware.Token).Claims.(jwt.MapClaims)
			c.Locals("user_id", int(claims["user_id"].(float64)))
			return c.Next()
		},
	})
}
