package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/fadilmartias/dilz_code/apps/backend/app/client"
	"github.com/fadilmartias/dilz_code/apps/backend/app/repositories"
	"github.com/fadilmartias/dilz_code/apps/backend/app/usecases"
	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"
	"github.com/fadilmartias/dilz_code/apps/backend/config"
	"github.com/go-redis/redis/v8"
	"github.com/pquerna/otp/totp"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type AuthOTPService struct {
	Redis          *redis.Client
	UserRepository *repositories.UserRepository
}

func NewAuthOTPService(redisClient *redis.Client, userRepository *repositories.UserRepository) *AuthOTPService {
	return &AuthOTPService{Redis: redisClient, UserRepository: userRepository}
}

func (s *AuthOTPService) Register2FA(email string) (string, string, error) {
	appConfig := config.LoadAppConfig()
	var issuer string
	if appConfig.Env != "production" {
		issuer = appConfig.Name + " " + cases.Title(language.English, cases.NoLower).String(appConfig.Env)
	} else {
		issuer = appConfig.Name
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: email,
	})
	if err != nil {
		return "", "", err
	}

	if err := s.Redis.Set(context.Background(), s.totpKey(email), key.Secret(), 5*time.Minute).Err(); err != nil {
		return "", "", err
	}

	return key.Secret(), key.URL(), nil
}

func (s *AuthOTPService) VerifyUserTOTPSecret(email string, code string, isLogin bool) error {
	var secret string
	if !isLogin {
		secretRedis, err := s.Redis.Get(context.Background(), s.totpKey(email)).Result()
		if err != nil {
			return err
		}
		secret = secretRedis
	} else {
		userDB, err := s.UserRepository.FindByEmail(email)
		if err != nil {
			return err
		}
		if userDB.TOTPSecret == nil {
			return errors.New("totp_secret not found")
		}
		secret = *userDB.TOTPSecret
	}

	if secret == "" {
		return errors.New("totp_secret not found")
	}

	if !totp.Validate(code, secret) {
		return errors.New("kode TOTP tidak valid")
	}

	if err := s.Redis.Del(context.Background(), s.totpKey(email)).Err(); err != nil {
		return err
	}

	userDB, err := s.UserRepository.FindByEmail(email)
	if err != nil {
		return err
	}
	if userDB.TOTPSecret == nil {
		userDB.TOTPSecret = &secret
		_, _ = s.UserRepository.UpdateWithModel(userDB, map[string]any{"totp_secret": secret})
	}
	return nil
}

func (s *AuthOTPService) Disable2FA(id string) error {
	_, err := s.UserRepository.UpdateByID(id, map[string]any{"totp_secret": nil})
	return err
}

func (s *AuthOTPService) SendOTP(ctx context.Context, targetType string, target string, subject string, name string) error {
	otp := strconv.Itoa(utils.GenerateRandomNumber(6))
	if err := s.Redis.Set(ctx, s.otpKey(subject, target), otp, 5*time.Minute).Err(); err != nil {
		return err
	}
	otpUsecase := usecases.NewSendOTPUsecase()
	otpUsecase.Execute(targetType, target, subject, otp, name)
	return nil
}

func (s *AuthOTPService) VerifyOTP(ctx context.Context, target string, subject string, otp string) error {
	storedOtp, err := s.Redis.Get(ctx, s.otpKey(subject, target)).Result()
	if err != nil {
		return err
	}
	if storedOtp != otp {
		return errors.New("kode OTP tidak valid")
	}
	return nil
}

func (s *AuthOTPService) totpKey(email string) string {
	return client.Key(fmt.Sprintf("totp:%s", email))
}

func (s *AuthOTPService) otpKey(subject, target string) string {
	return client.Key(fmt.Sprintf("otp:%s:%s", subject, target))
}
