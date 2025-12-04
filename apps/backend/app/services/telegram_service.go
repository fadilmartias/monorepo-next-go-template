package services

import (
	"github.com/fadilmartias/dilz_code/apps/backend/config"
	"github.com/go-resty/resty/v2"
)

type TelegramService struct {
}

func NewTelegramService() *TelegramService {
	return &TelegramService{}
}

func (s *TelegramService) Request(endpoint string, body map[string]any) (*resty.Response, error) {
	telegramConfig := config.LoadTelegramConfig()
	client := resty.New().
		SetBaseURL(telegramConfig.BaseURL+"/bot"+telegramConfig.BotToken).
		SetHeader("Content-Type", "application/json")
	return client.R().
		SetBody(body).
		Post(endpoint)
}

func (s *TelegramService) SendMessage(text string) (*resty.Response, error) {
	body := map[string]any{
		"chat_id": config.LoadTelegramConfig().DefaultChatID,
		"text":    text,
	}
	return s.Request("/sendMessage", body)
}
