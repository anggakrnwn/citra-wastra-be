package config

import (
	"citra-wastra-be/models"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {

	if err := godotenv.Load(); err != nil {
		fmt.Println(".env file not found!")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("gagal koneksi database", err)
	}

	db.Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto";`)

	fmt.Println("DB connected!")

	err = db.AutoMigrate(
		&models.User{},
		&models.UserProfile{},
		&models.Badge{},
		&models.UserBadge{},
		&models.Shop{},
		&models.BatikCatalog{},
	)

	if err != nil {
		log.Fatal("gagal migrasi database:", err)
	}

	fmt.Println("DB migrasi sukses!")

	return db
}
