package controllers_v1

import (
	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"
	"github.com/go-redis/redis/v8"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type CategoryController struct {
	BaseController
	DB    *gorm.DB
	Redis *redis.Client
}

func NewCategoryController(db *gorm.DB, redis *redis.Client) *CategoryController {
	return &CategoryController{DB: db, Redis: redis}
}

func (ctrl *CategoryController) Index(c fiber.Ctx) error {
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Berhasil mendapatkan data category",
		Data:    nil,
	})
}

func (ctrl *CategoryController) Show(c fiber.Ctx) error {
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Berhasil mendapatkan data category",
		Data:    nil,
	})
}

func (ctrl *CategoryController) Store(c fiber.Ctx) error {
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Berhasil menambahkan category",
		Data:    nil,
	})
}

func (ctrl *CategoryController) Update(c fiber.Ctx) error {
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Berhasil mengupdate category",
		Data:    nil,
	})
}

func (ctrl *CategoryController) Destroy(c fiber.Ctx) error {
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Berhasil menghapus category",
		Data:    nil,
	})
}
