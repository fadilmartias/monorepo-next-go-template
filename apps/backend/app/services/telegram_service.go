package services

import (
	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"
	"github.com/fadilmartias/dilz_code/apps/backend/config"
	"github.com/gofiber/fiber/v3/client"
)

type TelegramService struct {
}

func NewTelegramService() *TelegramService {
	return &TelegramService{}
}

func (s *TelegramService) Request(endpoint string, body map[string]any) (*client.Response, error) {
	telegramConfig := config.LoadTelegramConfig()
	return utils.Http().
		WithHeader("Content-Type", "application/json").
		WithJSON(body).
		Post(telegramConfig.BaseURL + "/bot" + telegramConfig.BotToken + endpoint)
}

func (s *TelegramService) SendMessage(text string) (*client.Response, error) {
	body := map[string]any{
		"chat_id": config.LoadTelegramConfig().DefaultChatID,
		"text":    text,
	}
	return s.Request("/sendMessage", body)
}
