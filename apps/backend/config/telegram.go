package config

import (
	"os"
	"sync"
)

type TelegramConfig struct {
	BotToken      string
	BaseURL       string
	DefaultChatID string
}

var (
	telegramConfig *TelegramConfig
	telegramOnce   sync.Once
)

func LoadTelegramConfig() *TelegramConfig {
	telegramOnce.Do(func() {
		telegramConfig = &TelegramConfig{
			BotToken:      os.Getenv("TELEGRAM_BOT_TOKEN"),
			BaseURL:       "https://api.telegram.org",
			DefaultChatID: "888560906",
		}
	})
	return telegramConfig
}
