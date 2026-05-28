package client

import (
	"fmt"
	"time"

	"github.com/fadilmartias/dilz_code/apps/backend/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

func ConnectDB() (*gorm.DB, error) {
	dbConfig := config.LoadDBConfig()
	appConfig := config.LoadAppConfig()

	var dialector gorm.Dialector

	// Tentukan Dialector berdasarkan driver yang dipilih
	switch dbConfig.Driver {
	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			dbConfig.Host,
			dbConfig.Port,
			dbConfig.User,
			dbConfig.Password,
			dbConfig.Name,
			dbConfig.SSLMode,
		)
		dialector = postgres.Open(dsn)

	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			dbConfig.User,
			dbConfig.Password,
			dbConfig.Host,
			dbConfig.Port,
			dbConfig.Name,
		)
		dialector = mysql.Open(dsn)

	default:
		return nil, fmt.Errorf("unsupported database driver: %s", dbConfig.Driver)
	}

	// Setup Logger
	var gormLogger glogger.Interface
	if appConfig.Env != "production" {
		gormLogger = glogger.Default.LogMode(glogger.Info)
	} else {
		gormLogger = glogger.Default.LogMode(glogger.Error)
	}

	// Buka koneksi menggunakan dialector yang sudah dinamis
	db, err := gorm.Open(dialector, &gorm.Config{Logger: gormLogger})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Setup Connection Pool
	if appConfig.Env != "production" {
		sqlDB.SetMaxIdleConns(5)
		sqlDB.SetMaxOpenConns(10)
		sqlDB.SetConnMaxLifetime(30 * time.Minute)
	} else {
		sqlDB.SetMaxIdleConns(20)
		sqlDB.SetMaxOpenConns(200)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	return db, nil
}
