package db

import (
	"fmt"
	"log"
	"piggybackend/config"
	"piggybackend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	config.LoadEnv()

	host := config.GetEnv("DB_HOST", "localhost")
	port := config.GetEnv("DB_PORT", "5433")
	user := config.GetEnv("DB_USER", "piggypocky")
	password := config.GetEnv("DB_PASSWORD", "piggy12123")
	dbname := config.GetEnv("DB_NAME", "piggydb")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// Auto-migrate tables
	DB.AutoMigrate(&models.User{}, &models.Amount{}, &models.Transaction{})
}
