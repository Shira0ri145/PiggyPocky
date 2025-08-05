package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username     string `gorm:"unique"`
	PasswordHash string
	Role         string
}

type LoginInput struct {
	Username string `json:"username" example:"admin"`
	Password string `json:"password" example:"admin123"`
}

type RegisterInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"` // admin or user
}

func FindUserByUsername(db *gorm.DB, username string) (*User, error) {
	var u User
	err := db.Where("username = ?", username).First(&u).Error
	return &u, err
}

func GetUserByID(db *gorm.DB, id uint) (*User, error) {
	var u User
	err := db.First(&u, id).Error
	return &u, err
}
