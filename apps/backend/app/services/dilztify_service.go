package services

import (
	"strings"

	"github.com/fadilmartias/dilz_code/apps/backend/app/client"
	"github.com/fadilmartias/dilz_code/apps/backend/app/requests"
	"github.com/fadilmartias/dilz_code/apps/backend/config"

	"github.com/gofiber/fiber/v3"
	fiberclient "github.com/gofiber/fiber/v3/client"
)

type DilztifyService struct {
}

func NewDilztifyService() *DilztifyService {
	return &DilztifyService{}
}

func (s *DilztifyService) DilztifySendMessage(c fiber.Ctx, body requests.DilztifySendMessageRequest) (*fiberclient.Response, error) {
	dilztifyConfig := config.LoadDilztifyConfig()
	body.Phone = formatPhoneNumber(body.Phone, "62")

	fullURL := dilztifyConfig.BaseURL + "/send/message"

	return client.HttpWithCtx(c).
		WithAuthToken(dilztifyConfig.APIKey).
		WithHeader("X-Device-ID", dilztifyConfig.DeviceID).
		WithJSON(body).
		Post(fullURL)
}

func formatPhoneNumber(phone string, countryCode string) string {
	phone = strings.TrimSpace(phone)
	phone = strings.TrimPrefix(phone, "+")
	if strings.HasPrefix(phone, countryCode) {
		return phone
	}
	phone = strings.TrimLeft(phone, "0")
	if !strings.HasPrefix(phone, countryCode) {
		phone = countryCode + phone
	}
	return phone
}
