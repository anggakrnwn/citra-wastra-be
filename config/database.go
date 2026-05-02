package config

import (
	"citra-wastra-be/config/database"
	"citra-wastra-be/models"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB() *gorm.DB {
	loadEnv()

	dsn := getDSN()

	config := &gorm.Config{
		PrepareStmt: false,
		Logger:      logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), config)
	if err != nil {
		log.Fatalf("CRITICAL: Gagal koneksi ke database: %v", err)
	}

	if err := setupDatabase(db); err != nil {
		log.Fatalf("CRITICAL: Gagal setup database: %v", err)
	}

	fmt.Println("database initialized successfully!")

	runSeeders(db)

	return db
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		if os.Getenv("APP_ENV") != "production" {
			log.Println("Info: .env file not found, using system environment variables")
		}
	}
}

func getDSN() string {
	dsn := os.Getenv("DATABASE_URL")
	if dsn != "" {
		if strings.Contains(dsn, "sslmode=verify-full") {
			return strings.Replace(dsn, "sslmode=verify-full", "sslmode=require", 1)
		}
		return dsn
	}

	sslmode := "disable"
	if os.Getenv("APP_ENV") == "production" {
		sslmode = "require"
	}

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		sslmode,
	)
}

func setupDatabase(db *gorm.DB) error {
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto";`).Error; err != nil {
		return fmt.Errorf("failed to create pgcrypto extension: %w", err)
	}

	allModels := []interface{}{
		&models.User{},
		&models.UserProfile{},
		&models.Badge{},
		&models.UserBadge{},
		&models.Shop{},
		&models.BatikCatalog{},
		&models.Batik{},
		&models.Island{},
		&models.Module{},
		&models.Level{},
		&models.UserLevelProgress{},
		&models.Question{},
		&models.AuditLog{},
		&models.SystemConfig{},
	}

	if err := db.AutoMigrate(allModels...); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	return nil
}

func runSeeders(db *gorm.DB) {
	if os.Getenv("APP_ENV") != "production" {
		database.SeedAll(db)
	} else {
		database.SeedAdmin(db)
		database.SeedConfigs(db)
	}
}
