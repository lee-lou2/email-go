package config

import (
	"gorm.io/gorm/logger"
	"log"
	"sync"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	dbInstance *gorm.DB
	dbOnce     sync.Once
)

// GetDB Database instance
func GetDB() *gorm.DB {
	dbOnce.Do(func() {
		dbFileName := "sqlite.db"
		db, err := gorm.Open(sqlite.Open(dbFileName), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Error),
		})
		if err != nil {
			log.Fatalf("failed to connect database: %v", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("failed to get generic database object: %v", err)
		}

		sqlDB.SetMaxOpenConns(10)
		sqlDB.SetMaxIdleConns(5)
		sqlDB.SetConnMaxLifetime(time.Hour)

		dbInstance = db
	})
	return dbInstance
}
