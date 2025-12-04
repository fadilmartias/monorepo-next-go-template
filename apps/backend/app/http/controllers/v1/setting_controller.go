package controllers_v1

import (
	"github.com/fadilmartias/dilz_code/apps/backend/app/models"
	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"
	"github.com/fadilmartias/dilz_code/apps/backend/config"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SettingController struct {
	BaseController
	DB    *gorm.DB
	Redis *config.RedisClient
}

func NewSettingController(db *gorm.DB, redis *config.RedisClient) *SettingController {
	return &SettingController{DB: db, Redis: redis}
}

func (ctrl *SettingController) Index(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Berhasil mendapatkan data setting",
		Data:    nil,
	})
}

func (ctrl *SettingController) Show(c *fiber.Ctx) error {
	key := c.Params("key")
	setting := models.Setting{}
	if err := ctrl.DB.Where("`key` = ?", key).First(&setting).Error; err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusInternalServerError,
			Message: "Gagal mendapatkan data setting",
		}, err)
	}
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Berhasil mendapatkan data setting",
		Data:    setting,
	})
}

func (ctrl *SettingController) Store(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Berhasil menambahkan setting",
		Data:    nil,
	})
}

func (ctrl *SettingController) Update(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Berhasil mengupdate setting",
		Data:    nil,
	})
}

func (ctrl *SettingController) Destroy(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Berhasil menghapus setting",
		Data:    nil,
	})
}
