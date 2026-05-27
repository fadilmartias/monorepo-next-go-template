package config

import (
	"os"
	"sync"
)

type DBConfig struct {
	Driver   string
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

var (
	dbConfig *DBConfig
	dbOnce   sync.Once
)

func LoadDBConfig() *DBConfig {
	dbOnce.Do(func() {
		if os.Getenv("DB_DRIVER") == "" {
			os.Setenv("DB_DRIVER", "mysql")
		}
		dbConfig = &DBConfig{
			Driver:   os.Getenv("DB_DRIVER"),
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("DB_SSLMODE"),
		}
	})
	return dbConfig
}
