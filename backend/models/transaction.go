// models/transaction.go
package models

import "time"

// Transaction DB model
type Transaction struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"column:user_id"`
	Amount    float64   `gorm:"column:amount"`
	Type      string    `gorm:"column:type"` // deposit / withdraw
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (Transaction) TableName() string {
	return "transactions"
}

// CoinInput is the request body for insert-coin
type CoinInput struct {
	UserID uint    `json:"user_id"`
	Amount float64 `json:"amount"`
}

// WithdrawInput is the request body for withdraw
type WithdrawInput struct {
	UserID uint    `json:"user_id"`
	Amount float64 `json:"amount"`
}
