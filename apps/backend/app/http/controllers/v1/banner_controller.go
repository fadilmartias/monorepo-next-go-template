package controllers_v1

import (
	"github.com/fadilmartias/dilz_code/apps/backend/app/models"
	"github.com/fadilmartias/dilz_code/apps/backend/app/services"
	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"

	"github.com/gofiber/fiber/v3"
)

type BannerController struct {
	BannerService *services.BannerService
}

func NewBannerController(bannerService *services.BannerService) *BannerController {
	return &BannerController{BannerService: bannerService}
}

// Index mengambil daftar banner yang aktif dan valid
// @Summary Ambil daftar banner aktif
// @Description Mengambil daftar banner yang berstatus aktif dan masih berada dalam masa berlaku (valid). Mendukung limitasi jumlah data.
// @Tags Banners
// @Accept json
// @Produce json
// @Param limit query int false "Batasan jumlah banner yang dikembalikan" default(10)
// @Success 200 {object} utils.OrderedSuccessResponse "Berhasil mendapatkan data banner"
// @Failure 500 {object} utils.OrderedErrorResponse "Gagal mendapatkan data banner"
// @Router /v1/banners [get]
func (ctrl *BannerController) Index(c fiber.Ctx) error {
	limit := fiber.Query[int](c, "limit", 10)
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

// Process melakukan operasi Insert atau Update pada banner
// @Summary Tambah atau Update Banner
// @Description Memproses data payload banner. Jika field ID pada payload kosong, maka akan melakukan Insert (Tambah Banner). Jika field ID terisi, akan melakukan Update pada banner tersebut.
// @Tags Banners
// @Accept json
// @Produce json
// @Param payload body models.Banner true "Data JSON Banner"
// @Success 200 {object} utils.OrderedSuccessResponse "Berhasil menambahkan atau mengupdate banner"
// @Failure 400 {object} utils.OrderedErrorResponse "Gagal memproses data banner dari body request"
// @Failure 500 {object} utils.OrderedErrorResponse "Gagal memproses/menyimpan banner ke database"
// @Router /v1/banners/process [post]
func (ctrl *BannerController) Process(c fiber.Ctx) error {
	banner := models.Banner{}
	if err := c.Bind().Body(&banner); err != nil {
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
