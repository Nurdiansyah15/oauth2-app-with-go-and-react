package database

import (
	"auth-server/models"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = db

	log.Println("Clearing existing data...")

	if err := DB.Exec("DELETE FROM users").Error; err != nil {
		log.Fatalf("Failed to clear users table: %v", err)
	}

	if err = DB.Exec("DELETE FROM clients").Error; err != nil {
		log.Fatalf("Failed to clear clients table: %v", err)
	}

	if err = DB.Exec("DELETE FROM auth_codes").Error; err != nil {
		log.Fatalf("Failed to clear auth_codes table: %v", err)
	}

	if err = DB.Exec("DELETE FROM sessions").Error; err != nil {
		log.Fatalf("Failed to clear sessions table: %v", err)
	}

	return DB.AutoMigrate(&models.User{}, &models.Client{}, &models.AuthCode{}, &models.Session{})
}
