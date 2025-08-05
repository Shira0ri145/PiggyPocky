package handlers

import (
	"piggybackend/db"
	"piggybackend/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// InsertCoin godoc
// @Summary Insert coin from ESP32
// @Description Accept coin and update user balance (ESP32 only)
// @Tags Transaction
// @Accept json
// @Produce json
// @Param X-ESP-KEY header string true "ESP API Key"
// @Param coin body models.CoinInput true "Coin data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/insert-coin [post]
func InsertCoin(c *fiber.Ctx) error {
	var input models.CoinInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}

	tx := models.Transaction{
		UserID: input.UserID,
		Amount: input.Amount,
		Type:   "deposit",
	}
	if err := db.DB.Create(&tx).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot insert transaction"})
	}

	if err := db.DB.Model(&models.Amount{}).
		Where("user_id = ?", input.UserID).
		Update("balance", gorm.Expr("balance + ?", input.Amount)).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot update balance"})
	}

	return c.JSON(fiber.Map{"message": "coin accepted"})
}

// WithdrawCoin godoc
// @Summary Withdraw money from user
// @Description Withdraw from balance if sufficient
// @Tags Transaction
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param withdraw body models.WithdrawInput true "Withdraw data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/withdraw [post]
func WithdrawCoin(c *fiber.Ctx) error {
	var input models.WithdrawInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}

	var userAmount models.Amount
	if err := db.DB.Where("user_id = ?", input.UserID).First(&userAmount).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	if userAmount.Balance < input.Amount {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "insufficient balance"})
	}

	tx := models.Transaction{
		UserID: input.UserID,
		Amount: -input.Amount,
		Type:   "withdraw",
	}
	if err := db.DB.Create(&tx).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot insert transaction"})
	}

	if err := db.DB.Model(&models.Amount{}).
		Where("user_id = ?", input.UserID).
		Update("balance", gorm.Expr("balance - ?", input.Amount)).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot update balance"})
	}

	return c.JSON(fiber.Map{"message": "withdraw successful"})
}

// GetTransactions godoc
// @Summary Get user transactions
// @Description Get list of transactions sorted by date and time (latest first)
// @Tags Transaction
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Transaction
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/transactions [get]
func GetTransactions(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(float64) // ดึงจาก JWT token

	var txs []models.Transaction
	if err := db.DB.
		Where("user_id = ?", uint(userID)).
		Order("created_at DESC").
		Find(&txs).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "cannot fetch transactions",
		})
	}

	return c.JSON(txs)
}

// GetMonthlySummary godoc
// @Summary Get monthly summary of transactions
// @Description Return summary of deposit and withdraw amounts grouped by month
// @Tags Transaction
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/transactions/monthly [get]
func GetMonthlySummary(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(float64)

	type Result struct {
		YearMonth     string  `json:"month"`      // "2025-07"
		TotalDeposit  float64 `json:"deposit"`    // sum of deposit
		TotalWithdraw float64 `json:"withdraw"`   // sum of withdraw
		NetAmount     float64 `json:"net_amount"` // deposit - withdraw
	}

	var results []Result
	err := db.DB.Raw(`
		SELECT
			TO_CHAR(created_at, 'YYYY-MM') AS year_month,
			COALESCE(SUM(CASE WHEN type = 'deposit' THEN amount ELSE 0 END), 0) AS total_deposit,
			COALESCE(SUM(CASE WHEN type = 'withdraw' THEN -amount ELSE 0 END), 0) AS total_withdraw,
			COALESCE(SUM(amount), 0) AS net_amount
		FROM transactions
		WHERE user_id = ?
		GROUP BY year_month
		ORDER BY year_month DESC
	`, uint(userID)).Scan(&results).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "cannot fetch monthly summary",
		})
	}

	return c.JSON(results)
}
