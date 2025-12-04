package services

import (
	"github.com/fadilmartias/dilz_code/apps/backend/app/requests"
	"github.com/fadilmartias/dilz_code/apps/backend/config"
	"github.com/go-resty/resty/v2"
)

type FonnteService struct {
}

func NewFonnteService() *FonnteService {
	return &FonnteService{}
}

func (s *FonnteService) FonnteSendMessage(body requests.FonnteSendMessageRequest) (*resty.Response, error) {
	fonnteConfig := config.LoadFonnteConfig()
	client := resty.New().
		SetBaseURL(fonnteConfig.BaseURL).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fonnteConfig.Token)
	return client.R().
		SetBody(body).
		Post("/send")
}
