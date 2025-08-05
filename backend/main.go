package main

import (
	"piggybackend/db"
	"piggybackend/handlers"
	"piggybackend/middleware"

	_ "piggybackend/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
)

// @title PiggyPocky API
// @version 1.0
// @description Backend API for PiggyPocky
// @host localhost:3000
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token. Example: "Bearer xxxxxx"

func main() {
	db.InitDB()
	app := fiber.New()
	app.Use(logger.New())
	app.Use(middleware.Cors())
	db.InitRedis()

	// Swagger + Auth
	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Post("/login", handlers.Login)

	// ✅ ต้องมาก่อน group /api
	app.Post("/api/insert-coin", middleware.OnlyESP(), handlers.InsertCoin)

	// ✅ กลุ่ม API ที่ต้องใช้ JWT
	api := app.Group("/api", middleware.Protect())
	api.Get("/profile", handlers.Profile)
	api.Get("/amount", handlers.Amount)
	api.Post("/register", middleware.OnlyAdmin(), handlers.Register)
	api.Post("/logout", handlers.Logout)
	api.Post("/withdraw", handlers.WithdrawCoin)
	api.Get("/transactions", handlers.GetTransactions)
	api.Get("/transactions/monthly", handlers.GetMonthlySummary)

	app.Listen(":3000")
}
