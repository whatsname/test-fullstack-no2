package database

import (
	"log"

	"backend/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(dbPath string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database (%s): %v", dbPath, err)
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	return db
}
