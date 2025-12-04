package job_handlers

import (
	job_tasks "github.com/fadilmartias/dilz_code/apps/backend/app/jobs/tasks"
	"github.com/fadilmartias/dilz_code/apps/backend/app/services"
	"github.com/fadilmartias/dilz_code/apps/backend/config"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

func NewHandler(db *gorm.DB, redis *config.RedisClient) *asynq.ServeMux {
	mux := asynq.NewServeMux()
	telegramService := services.NewTelegramService()

	// EMAIL HANDLER
	emailHandler := NewEmailHandler(telegramService)
	mux.HandleFunc(job_tasks.TypeEmailVerification, emailHandler.HandleVerification)
	mux.HandleFunc(job_tasks.TypeEmailResetPassword, emailHandler.HandleResetPassword)

	// OTP HANDLER
	otpHandler := NewOTPHandler(db, telegramService, services.NewFonnteService())
	mux.HandleFunc(job_tasks.TypeSendOTP, otpHandler.SendOTP)

	// PRODUCT ORDER HANDLER
	return mux
}
