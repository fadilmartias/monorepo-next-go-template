package services

import "github.com/fadilmartias/dilz_code/apps/backend/app/usecases"

type EmailService struct {
	SendEmailVerificationUsecase  *usecases.SendEmailVerificationUsecase
	SendResetPasswordEmailUsecase *usecases.SendResetPasswordEmailUsecase
}

func NewEmailService(sendEmailVerificationUsecase *usecases.SendEmailVerificationUsecase, sendResetPasswordEmailUsecase *usecases.SendResetPasswordEmailUsecase) *EmailService {
	return &EmailService{SendEmailVerificationUsecase: sendEmailVerificationUsecase, SendResetPasswordEmailUsecase: sendResetPasswordEmailUsecase}
}

func (s *EmailService) SendEmailVerification(email string, token string, name string) error {
	return s.SendEmailVerificationUsecase.Execute(email, token, name)
}

func (s *EmailService) SendEmailResetPassword(email string, token string, name string) error {
	return s.SendResetPasswordEmailUsecase.Execute(email, token, name)
}
