package controllers_v1

import (
	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"
	"github.com/fadilmartias/dilz_code/apps/backend/config"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type PaymentGatewayController struct {
	BaseController
	DB    *gorm.DB
	Redis *config.RedisClient
}

func NewPaymentGatewayController(db *gorm.DB, redis *config.RedisClient) *PaymentGatewayController {
	return &PaymentGatewayController{DB: db, Redis: redis}
}

func (ctrl *PaymentGatewayController) Index(c fiber.Ctx) error {
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Berhasil mendapatkan data payment gateway",
		Data:    nil,
	})
}

func (ctrl *PaymentGatewayController) Show(c fiber.Ctx) error {
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Berhasil mendapatkan data payment gateway",
		Data:    nil,
	})
}

func (ctrl *PaymentGatewayController) Store(c fiber.Ctx) error {
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Berhasil menambahkan payment gateway",
		Data:    nil,
	})
}

func (ctrl *PaymentGatewayController) Update(c fiber.Ctx) error {
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Berhasil mengupdate payment gateway",
		Data:    nil,
	})
}

func (ctrl *PaymentGatewayController) Destroy(c fiber.Ctx) error {
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Berhasil menghapus payment gateway",
		Data:    nil,
	})
}
