package controllers_v1

import (
	"github.com/fadilmartias/dilz_code/apps/backend/app/models"
	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"
	"github.com/go-redis/redis/v8"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type SettingController struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func NewSettingController(db *gorm.DB, redis *redis.Client) *SettingController {
	return &SettingController{DB: db, Redis: redis}
}

// Index mengambil daftar semua pengaturan
// @Summary Ambil daftar pengaturan
// @Description Mengambil semua data konfigurasi dan pengaturan aplikasi.
// @Tags Settings
// @Accept json
// @Produce json
// @Success 200 {object} utils.OrderedSuccessResponse "Berhasil mendapatkan data setting"
// @Router /v1/settings [get]
func (ctrl *SettingController) Index(c fiber.Ctx) error {
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Berhasil mendapatkan data setting",
		Data:    nil,
	})
}

// Show mengambil detail pengaturan berdasarkan key
// @Summary Ambil detail pengaturan
// @Description Mengambil nilai konfigurasi atau pengaturan spesifik berdasarkan 'key' (kunci) pengaturannya.
// @Tags Settings
// @Accept json
// @Produce json
// @Param key path string true "Key string dari pengaturan (contoh: 'site_title')"
// @Success 200 {object} utils.OrderedSuccessResponse "Berhasil mendapatkan data setting"
// @Failure 500 {object} utils.OrderedErrorResponse "Gagal mendapatkan data setting"
// @Router /v1/settings/{key} [get]
func (ctrl *SettingController) Show(c fiber.Ctx) error {
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

// Store menambahkan data pengaturan baru
// @Summary Tambah pengaturan baru
// @Description Membuat data konfigurasi pengaturan baru ke dalam database.
// @Tags Settings
// @Accept json
// @Produce json
// @Param payload body map[string]interface{} true "Data JSON pengaturan baru"
// @Success 201 {object} utils.OrderedSuccessResponse "Berhasil menambahkan setting"
// @Router /v1/settings [post]
func (ctrl *SettingController) Store(c fiber.Ctx) error {
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Berhasil menambahkan setting",
		Data:    nil,
	})
}

// Update memodifikasi data pengaturan yang sudah ada
// @Summary Update pengaturan
// @Description Memodifikasi data pengaturan berdasarkan ID.
// @Tags Settings
// @Accept json
// @Produce json
// @Param id path string true "ID pengaturan"
// @Param payload body map[string]interface{} true "Data JSON pembaruan pengaturan"
// @Success 200 {object} utils.OrderedSuccessResponse "Berhasil mengupdate setting"
// @Router /v1/settings/{id} [put]
func (ctrl *SettingController) Update(c fiber.Ctx) error {
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Berhasil mengupdate setting",
		Data:    nil,
	})
}

// Destroy menghapus data pengaturan
// @Summary Hapus pengaturan
// @Description Menghapus data pengaturan dari database berdasarkan ID.
// @Tags Settings
// @Accept json
// @Produce json
// @Param id path string true "ID pengaturan"
// @Success 200 {object} utils.OrderedSuccessResponse "Berhasil menghapus setting"
// @Router /v1/settings/{id} [delete]
func (ctrl *SettingController) Destroy(c fiber.Ctx) error {
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Berhasil menghapus setting",
		Data:    nil,
	})
}
