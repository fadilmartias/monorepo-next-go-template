package controllers_v1

import (
	"github.com/fadilmartias/dilz_code/apps/backend/app/services"
	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"

	"github.com/gofiber/fiber/v2"
)

type ActivityLogController struct {
	ActivityLogService *services.ActivityLogService
}

func NewActivityLogController(activityLogService *services.ActivityLogService) *ActivityLogController {
	return &ActivityLogController{ActivityLogService: activityLogService}
}
func (ctrl *ActivityLogController) Show(c *fiber.Ctx) error {
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
