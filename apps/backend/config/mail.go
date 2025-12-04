package config

import (
	"log"
	"os"
	"strconv"
	"sync"
)

type MailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	AppURL   string
}

var (
	mailConfig *MailConfig
	mailOnce   sync.Once
)

func LoadMailConfig() *MailConfig {
	mailOnce.Do(func() {
		port, err := strconv.Atoi(os.Getenv("MAIL_PORT"))
		if err != nil {
			log.Fatalf("MAIL_PORT harus berupa angka: %v", err)
		}
		mailConfig = &MailConfig{
			Host:     os.Getenv("MAIL_HOST"),
			Port:     port,
			Username: os.Getenv("MAIL_USERNAME"),
			Password: os.Getenv("MAIL_PASSWORD"),
			From:     os.Getenv("MAIL_FROM"),
			AppURL:   os.Getenv("APP_URL"),
		}
	})
	return mailConfig
}
