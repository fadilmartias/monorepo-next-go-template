package job_handlers

import (
	"context"
	"fmt"

	"github.com/bytedance/sonic"
	job_tasks "github.com/fadilmartias/dilz_code/apps/backend/app/jobs/tasks"
	"github.com/fadilmartias/dilz_code/apps/backend/app/logger"
	"github.com/fadilmartias/dilz_code/apps/backend/app/requests"
	"github.com/fadilmartias/dilz_code/apps/backend/app/services"
	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

type OTPHandler struct {
	DB              *gorm.DB
	TelegramService *services.TelegramService
	FonnteService   *services.FonnteService
}

func NewOTPHandler(db *gorm.DB, telegramService *services.TelegramService, fonnteService *services.FonnteService) *OTPHandler {
	return &OTPHandler{
		DB:              db,
		TelegramService: telegramService,
		FonnteService:   fonnteService,
	}
}

func (h *OTPHandler) SendOTP(ctx context.Context, t *asynq.Task) error {
	var payload job_tasks.SendOTPPayload
	if err := sonic.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	// Retry check
	retries, _ := asynq.GetRetryCount(ctx)
	maxRetry, _ := asynq.GetMaxRetry(ctx)
	if retries == maxRetry {
		h.TelegramService.SendMessage("Worker: Gagal mengirim OTP: " + string(t.Payload()))
		return fmt.Errorf("max retry reached")
	}

	switch payload.TargetType {
	case "email":
		emailHandler := NewEmailHandler(h.TelegramService)
		return emailHandler.HandleOTP(ctx, t)
	case "wa":
		message := "Kode OTP kamu untuk DilZ Topup: " + payload.OTP + "\nKode ini berlaku selama 5 menit. Jangan bagikan kode ini kepada siapa pun, termasuk pihak yang mengatasnamakan DilZ Topup."
		resp, err := h.FonnteService.FonnteSendMessage(requests.FonnteSendMessageRequest{
			Target:  utils.StringPtr(payload.Target),
			Message: utils.StringPtr(message),
		})
		if err != nil {
			return fmt.Errorf("gagal mengirim pesan: %w", err)
		}
		logger.Infof("Pesan berhasil dikirim: %v", resp)
		return nil
	default:
		return fmt.Errorf("unsupported target type: %s", payload.TargetType)
	}
}
