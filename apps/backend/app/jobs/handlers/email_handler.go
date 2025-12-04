package job_handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"time"

	job_tasks "github.com/fadilmartias/dilz_code/apps/backend/app/jobs/tasks"
	"github.com/fadilmartias/dilz_code/apps/backend/app/services"
	"github.com/fadilmartias/dilz_code/apps/backend/config"
	"github.com/hibiken/asynq"
	"gopkg.in/gomail.v2"
)

type BaseEmailData struct {
	APP_NAME          string
	LOGO_URL          string
	COMPANY_NAME      string
	USER_NAME         string
	WHATSAPP_NUMBER   string
	TELEGRAM_USERNAME string
	FACEBOOK_URL      string
	YOUTUBE_URL       string
	TWITTER_URL       string
	INSTAGRAM_URL     string
	COMPANY_ADDRESS   string
	SUPPORT_EMAIL     string
	WEBSITE_URL       string
	YEAR              string
	USER_EMAIL        string
}

type ResetPasswordData struct {
	BaseEmailData
	ICON_LOCK_URL  string
	RESET_URL      string
	EXPIRY_MINUTES int
}

type VerificationEmailData struct {
	BaseEmailData
	ICON_MAIL_URL    string
	VERIFICATION_URL string
	EXPIRY_MINUTES   int
}

type OTPData struct {
	BaseEmailData
	ICON_LOCK_URL  string
	OTP            string
	EXPIRY_MINUTES int
}

type EmailHandler struct {
	TelegramService *services.TelegramService
	AppConfig       *config.AppConfig
	MailConfig      *config.MailConfig
}

func NewEmailHandler(telegramService *services.TelegramService) *EmailHandler {
	return &EmailHandler{TelegramService: telegramService, AppConfig: config.LoadAppConfig(), MailConfig: config.LoadMailConfig()}
}

func (s *EmailHandler) BuildBaseEmailData(name, email string) BaseEmailData {
	return BaseEmailData{
		APP_NAME:          s.AppConfig.Name,
		LOGO_URL:          "https://dilztopup.com/assets/images/logos/logo-dark.png",
		USER_NAME:         name,
		WHATSAPP_NUMBER:   s.AppConfig.WhatsappNumber,
		TELEGRAM_USERNAME: s.AppConfig.TelegramUsername,
		COMPANY_NAME:      s.AppConfig.CompanyName,
		FACEBOOK_URL:      "https://facebook.com/dilz_topup",
		YOUTUBE_URL:       "https://youtube.com/dilz_topup",
		TWITTER_URL:       "https://twitter.com/dilz_topup",
		INSTAGRAM_URL:     "https://instagram.com/dilz_topup",
		COMPANY_ADDRESS:   "Riau, Indonesia",
		SUPPORT_EMAIL:     "support@dilztopup.com",
		WEBSITE_URL:       "https://dilztopup.com",
		YEAR:              time.Now().Format("2006"),
		USER_EMAIL:        email,
	}
}

func (h *EmailHandler) HandleResetPassword(ctx context.Context, t *asynq.Task) error {
	var payload job_tasks.EmailResetPasswordPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	if h.isMaxRetry(ctx, t, "reset password") {
		return fmt.Errorf("max retry reached")
	}

	tmpl, err := template.ParseFiles("templates/email/reset_password.html")
	if err != nil {
		return err
	}

	baseEmailData := h.BuildBaseEmailData(payload.Name, payload.To)

	data := ResetPasswordData{
		BaseEmailData:  baseEmailData,
		ICON_LOCK_URL:  "https://cdn-icons-png.flaticon.com/512/3039/3039427.png",
		RESET_URL:      h.AppConfig.FEURL + "/reset-password?token=" + payload.Token,
		EXPIRY_MINUTES: 5,
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}

	return h.sendMail(payload.To, "Reset Password Request", body.String())
}

func (h *EmailHandler) HandleVerification(ctx context.Context, t *asynq.Task) error {
	var payload job_tasks.EmailVerificationPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	if h.isMaxRetry(ctx, t, "verifikasi email") {
		return fmt.Errorf("max retry reached")
	}

	appConfig := config.LoadAppConfig()
	tmpl, err := template.ParseFiles("templates/email/verification_email.html")
	if err != nil {
		return err
	}

	baseEmailData := h.BuildBaseEmailData(payload.Name, payload.To)

	data := VerificationEmailData{
		BaseEmailData:    baseEmailData,
		ICON_MAIL_URL:    "https://cdn-icons-png.flaticon.com/512/2099/2099131.png",
		VERIFICATION_URL: appConfig.FEURL + "/verify-email?token=" + payload.Token,
		EXPIRY_MINUTES:   5,
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}

	return h.sendMail(payload.To, "Email Verification", body.String())
}

func (h *EmailHandler) HandleOTP(ctx context.Context, t *asynq.Task) error {
	var payload job_tasks.SendOTPPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	if h.isMaxRetry(ctx, t, "OTP") {
		return fmt.Errorf("max retry reached")
	}

	tmpl, err := template.ParseFiles("templates/email/otp.html")
	if err != nil {
		return err
	}

	baseEmailData := h.BuildBaseEmailData(payload.Name, payload.Target)

	data := OTPData{
		ICON_LOCK_URL:  "https://cdn-icons-png.flaticon.com/512/3039/3039427.png",
		BaseEmailData:  baseEmailData,
		OTP:            payload.OTP,
		EXPIRY_MINUTES: 5,
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}

	return h.sendMail(payload.Target, payload.Subject, body.String())
}

// helper function
func (h *EmailHandler) sendMail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", h.MailConfig.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(h.MailConfig.Host, h.MailConfig.Port, h.MailConfig.Username, h.MailConfig.Password)

	if err := d.DialAndSend(m); err != nil {
		log.Println("Error sending email:", err)
		return err
	}
	return nil
}

func (h *EmailHandler) isMaxRetry(ctx context.Context, t *asynq.Task, taskName string) bool {
	retries, ok := asynq.GetRetryCount(ctx)
	if !ok {
		return false
	}
	maxRetry, ok := asynq.GetMaxRetry(ctx)
	if !ok {
		return false
	}
	if retries == maxRetry {
		h.TelegramService.SendMessage("Worker: Gagal mengirim " + taskName + ": " + string(t.Payload()))
		return true
	}
	return false
}
