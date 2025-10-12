package db

import (
	"NotifyX/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

var DB *gorm.DB

func ConnectDb() {
	dsn := "host=localhost user=postgres password=password dbname=testing sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	// ✅ Get underlying sql.DB for connection pooling
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("failed to get sql.DB from gorm: ", err)
	}

	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute) // Recycle connection after 30 mins
	sqlDB.SetConnMaxIdleTime(10 * time.Minute) // Idle timeout

	// ✅ Auto migrate models
	err = db.AutoMigrate(&models.Notification{}, &models.NotificationRecipient{}, &models.ProviderLog{})
	if err != nil {
		log.Fatal("AutoMigrate error: ", err)
	}

	DB = db
}
