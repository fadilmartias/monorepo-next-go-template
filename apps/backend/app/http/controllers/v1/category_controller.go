package controllers_v1

import (
	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"
	"github.com/go-redis/redis/v8"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type CategoryController struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func NewCategoryController(db *gorm.DB, redis *redis.Client) *CategoryController {
	return &CategoryController{DB: db, Redis: redis}
}

// Index mengambil semua data kategori
// @Summary Ambil daftar kategori
// @Description Mengambil semua data kategori yang tersedia.
// @Tags Categories
// @Accept json
// @Produce json
// @Success 200 {object} utils.OrderedSuccessResponse "Berhasil mendapatkan data category"
// @Router /v1/categories [get]
func (ctrl *CategoryController) Index(c fiber.Ctx) error {
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Berhasil mendapatkan data category",
		Data:    nil,
	})
}

// Show mengambil detail data kategori berdasarkan ID
// @Summary Ambil detail kategori
// @Description Mengambil satu baris data kategori secara spesifik menggunakan parameter ID.
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path string true "ID kategori"
// @Success 200 {object} utils.OrderedSuccessResponse "Berhasil mendapatkan data category"
// @Router /v1/categories/{id} [get]
func (ctrl *CategoryController) Show(c fiber.Ctx) error {
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Berhasil mendapatkan data category",
		Data:    nil,
	})
}

// Store menambahkan data kategori baru
// @Summary Tambah kategori baru
// @Description Membuat entri data kategori baru ke dalam database.
// @Tags Categories
// @Accept json
// @Produce json
// @Param payload body map[string]interface{} true "Data JSON kategori baru"
// @Success 201 {object} utils.OrderedSuccessResponse "Berhasil menambahkan category"
// @Router /v1/categories [post]
func (ctrl *CategoryController) Store(c fiber.Ctx) error {
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Berhasil menambahkan category",
		Data:    nil,
	})
}

// Update memodifikasi data kategori yang sudah ada
// @Summary Update kategori
// @Description Memodifikasi keseluruhan data kategori berdasarkan ID.
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path string true "ID kategori"
// @Param payload body map[string]interface{} true "Data JSON pembaruan kategori"
// @Success 200 {object} utils.OrderedSuccessResponse "Berhasil mengupdate category"
// @Router /v1/categories/{id} [put]
func (ctrl *CategoryController) Update(c fiber.Ctx) error {
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Berhasil mengupdate category",
		Data:    nil,
	})
}

// Destroy menghapus data kategori
// @Summary Hapus kategori
// @Description Menghapus data kategori dari database berdasarkan ID.
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path string true "ID kategori"
// @Success 200 {object} utils.OrderedSuccessResponse "Berhasil menghapus category"
// @Router /v1/categories/{id} [delete]
func (ctrl *CategoryController) Destroy(c fiber.Ctx) error {
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Berhasil menghapus category",
		Data:    nil,
	})
}