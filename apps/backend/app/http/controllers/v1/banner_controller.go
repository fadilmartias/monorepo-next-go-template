package controllers_v1

import (
	"github.com/fadilmartias/dilz_code/apps/backend/app/models"
	"github.com/fadilmartias/dilz_code/apps/backend/app/services"
	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"

	"github.com/gofiber/fiber/v2"
)

type BannerController struct {
	BannerService *services.BannerService
}

func NewBannerController(bannerService *services.BannerService) *BannerController {
	return &BannerController{BannerService: bannerService}
}

func (ctrl *BannerController) Index(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)
	banners, err := ctrl.BannerService.GetActiveAndValidBanners(limit)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusInternalServerError,
			Message: "Gagal mendapatkan data banner",
		}, err)
	}
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Berhasil mendapatkan data banner",
		Data:    banners,
	})
}

func (ctrl *BannerController) Process(c *fiber.Ctx) error {
	banner := models.Banner{}
	if err := c.BodyParser(&banner); err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusBadRequest,
			Message: "Gagal memproses data banner",
		}, err)
	}

	banner.TenantID = "1"
	result, err := ctrl.BannerService.Process(&banner)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusInternalServerError,
			Message: "Gagal memproses banner",
		}, err)
	}

	msg := "Berhasil menambahkan banner"
	if banner.ID != "" {
		msg = "Berhasil mengupdate banner"
	}

	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: msg,
		Data:    result,
	})
}
