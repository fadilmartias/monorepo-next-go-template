package controllers_v1

import (
	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"
	"github.com/go-redis/redis/v8"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type DashboardController struct {
	BaseController
	DB    *gorm.DB
	Redis *redis.Client
}

func NewDashboardController(db *gorm.DB, redis *redis.Client) *DashboardController {
	return &DashboardController{DB: db, Redis: redis}
}

func (ctrl *DashboardController) Index(c fiber.Ctx) error {

	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Berhasil mendapatkan data dashboard",
	})
}
