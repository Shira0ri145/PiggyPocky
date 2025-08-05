package models

import "gorm.io/gorm"

type Amount struct {
	ID      uint    `gorm:"primaryKey"`
	UserID  uint    `gorm:"column:user_id"`
	Balance float64 `gorm:"column:balance"`
}

type AmountResponse struct {
	UserID  uint    `json:"user_id"`
	Balance float64 `json:"balance"`
}

func (Amount) TableName() string {
	return "amount"
}

func GetAmountByUserID(db *gorm.DB, userID uint) (*Amount, error) {
	var a Amount
	err := db.Where("user_id = ?", userID).First(&a).Error
	return &a, err
}
