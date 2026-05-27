package controllers_v1

import (
	"github.com/fadilmartias/dilz_code/apps/backend/app/services"
	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"

	"github.com/gofiber/fiber/v3"
)

type ActivityLogController struct {
	ActivityLogService *services.ActivityLogService
}

func NewActivityLogController(activityLogService *services.ActivityLogService) *ActivityLogController {
	return &ActivityLogController{ActivityLogService: activityLogService}
}

// Show mengambil data Activity Log spesifik
// @Summary Ambil detail Activity Log
// @Description Mengambil satu baris data riwayat aktivitas (Activity Log) berdasarkan ID.
// @Tags Activity Log
// @Accept json
// @Produce json
// @Param id path string true "ID dari Activity Log"
// @Success 200 {object} utils.OrderedSuccessResponse "Berhasil mendapatkan data activity log"
// @Failure 404 {object} utils.OrderedErrorResponse "Activity Log tidak ditemukan"
// @Failure 500 {object} utils.OrderedErrorResponse "Terjadi kesalahan pada server"
// @Router /v1/activity-logs/{id} [get]
func (ctrl *ActivityLogController) Show(c fiber.Ctx) error {
	id := c.Params("id")
	activityLog, err := ctrl.ActivityLogService.FindByID(id)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusNotFound,
			Message: "Activity Log tidak ditemukan",
		}, err)
	}
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Berhasil mendapatkan data activity log",
		Data:    activityLog,
	})
}
