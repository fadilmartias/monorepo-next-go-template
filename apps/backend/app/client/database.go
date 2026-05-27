package client

import (
	"fmt"
	"time"

	"github.com/fadilmartias/dilz_code/apps/backend/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

func ConnectDB() (*gorm.DB, error) {
	dbConfig := config.LoadDBConfig()
	appConfig := config.LoadAppConfig()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name,
	)

	var gormLogger glogger.Interface
	if appConfig.Env != "production" {
		gormLogger = glogger.Default.LogMode(glogger.Info)
	} else {
		gormLogger = glogger.Default.LogMode(glogger.Error)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: gormLogger})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

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
