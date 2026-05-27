package controllers_v1

import (
	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"
	"github.com/go-redis/redis/v8"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type DashboardController struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func NewDashboardController(db *gorm.DB, redis *redis.Client) *DashboardController {
	return &DashboardController{DB: db, Redis: redis}
}

// Index mengambil data ringkasan untuk dashboard
// @Summary Ambil data dashboard
// @Description Mengambil ringkasan data statistik dan metrik utama untuk ditampilkan di halaman dashboard.
// @Tags Dashboard
// @Accept json
// @Produce json
// @Success 200 {object} utils.OrderedSuccessResponse "Berhasil mendapatkan data dashboard"
// @Router /v1/dashboard [get]
func (ctrl *DashboardController) Index(c fiber.Ctx) error {

	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Berhasil mendapatkan data dashboard",
	})
}
