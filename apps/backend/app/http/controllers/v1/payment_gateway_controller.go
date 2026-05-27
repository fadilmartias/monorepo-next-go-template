package controllers_v1

import (
	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"
	"github.com/go-redis/redis/v8"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type PaymentGatewayController struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func NewPaymentGatewayController(db *gorm.DB, redis *redis.Client) *PaymentGatewayController {
	return &PaymentGatewayController{DB: db, Redis: redis}
}

// Index mengambil daftar payment gateway
// @Summary Ambil daftar payment gateway
// @Description Mengambil semua data konfigurasi payment gateway yang tersedia di sistem.
// @Tags Payment Gateways
// @Accept json
// @Produce json
// @Success 200 {object} utils.OrderedSuccessResponse "Berhasil mendapatkan data payment gateway"
// @Router /v1/payment-gateways [get]
func (ctrl *PaymentGatewayController) Index(c fiber.Ctx) error {
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Berhasil mendapatkan data payment gateway",
		Data:    nil,
	})
}

// Show mengambil detail data payment gateway berdasarkan ID
// @Summary Ambil detail payment gateway
// @Description Mengambil satu baris data konfigurasi payment gateway secara spesifik menggunakan parameter ID.
// @Tags Payment Gateways
// @Accept json
// @Produce json
// @Param id path string true "ID Payment Gateway"
// @Success 200 {object} utils.OrderedSuccessResponse "Berhasil mendapatkan data payment gateway"
// @Router /v1/payment-gateways/{id} [get]
func (ctrl *PaymentGatewayController) Show(c fiber.Ctx) error {
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Berhasil mendapatkan data payment gateway",
		Data:    nil,
	})
}

// Store menambahkan data payment gateway baru
// @Summary Tambah payment gateway baru
// @Description Membuat entri data konfigurasi payment gateway baru ke dalam database.
// @Tags Payment Gateways
// @Accept json
// @Produce json
// @Param payload body map[string]interface{} true "Data JSON payment gateway baru"
// @Success 201 {object} utils.OrderedSuccessResponse "Berhasil menambahkan payment gateway"
// @Router /v1/payment-gateways [post]
func (ctrl *PaymentGatewayController) Store(c fiber.Ctx) error {
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Berhasil menambahkan payment gateway",
		Data:    nil,
	})
}

// Update memodifikasi data payment gateway yang sudah ada
// @Summary Update payment gateway
// @Description Memodifikasi keseluruhan data konfigurasi payment gateway berdasarkan ID.
// @Tags Payment Gateways
// @Accept json
// @Produce json
// @Param id path string true "ID Payment Gateway"
// @Param payload body map[string]interface{} true "Data JSON pembaruan payment gateway"
// @Success 200 {object} utils.OrderedSuccessResponse "Berhasil mengupdate payment gateway"
// @Router /v1/payment-gateways/{id} [put]
func (ctrl *PaymentGatewayController) Update(c fiber.Ctx) error {
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Berhasil mengupdate payment gateway",
		Data:    nil,
	})
}

// Destroy menghapus data payment gateway
// @Summary Hapus payment gateway
// @Description Menghapus data konfigurasi payment gateway dari database berdasarkan ID.
// @Tags Payment Gateways
// @Accept json
// @Produce json
// @Param id path string true "ID Payment Gateway"
// @Success 200 {object} utils.OrderedSuccessResponse "Berhasil menghapus payment gateway"
// @Router /v1/payment-gateways/{id} [delete]
func (ctrl *PaymentGatewayController) Destroy(c fiber.Ctx) error {
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Berhasil menghapus payment gateway",
		Data:    nil,
	})
}
