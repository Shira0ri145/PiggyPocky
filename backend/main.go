package main

import (
	"coin-backend/db"
	"coin-backend/handlers"
	"coin-backend/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	_ "coin-backend/docs"
)

// @title Coin API
// @version 1.0
// @description API for login, profile, amount and transactions
// @host localhost:3000
// @BasePath /
func main() {
	app := fiber.New()
	db.InitDB()

	app.Post("/login", handlers.Login)

	api := app.Group("/api", middleware.Protect())
	api.Get("/profile", handlers.Profile)
	api.Get("/amount", handlers.Amount)
	api.Get("/transactions", handlers.GetTransactions)
	api.Get("/transactions/:id", handlers.GetTransactionByID)

	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Listen(":3000")
}
