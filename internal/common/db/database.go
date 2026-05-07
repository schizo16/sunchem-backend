package db

import (
	"fmt"
	"log"

	"sunchem-backend/internal/common/config"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(cfg *config.Config) *gorm.DB {
	var dialector gorm.Dialector
	switch cfg.DBType {
	case "mysql":
		dialector = mysql.Open(cfg.DBDSN)
	default:
		dialector = sqlite.Open(cfg.DBDSN)
	}

	logLevel := logger.Info
	if cfg.AppEnv == "production" {
		logLevel = logger.Warn
	}

	db, err := gorm.Open(dialector, &gorm.Config{Logger: logger.Default.LogMode(logLevel)})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if cfg.DBType == "sqlite" {
		db.Exec("PRAGMA journal_mode=WAL")
		db.Exec("PRAGMA foreign_keys=ON")
	}

	fmt.Println("[DB] Connected successfully")
	return db
}
