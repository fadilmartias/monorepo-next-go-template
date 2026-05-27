package services

import (
	"github.com/fadilmartias/dilz_code/apps/backend/app/requests"
	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"
	"github.com/fadilmartias/dilz_code/apps/backend/config"
	"github.com/gofiber/fiber/v3/client"
)

type FonnteService struct {
}

func NewFonnteService() *FonnteService {
	return &FonnteService{}
}

func (s *FonnteService) FonnteSendMessage(body requests.FonnteSendMessageRequest) (*client.Response, error) {
	fonnteConfig := config.LoadFonnteConfig()
	return utils.Http().
		WithHeader("Authorization", fonnteConfig.Token).
		WithJSON(body).
		Post(fonnteConfig.BaseURL + "/send")
}
