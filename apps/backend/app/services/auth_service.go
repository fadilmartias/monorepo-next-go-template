package services

import (
	"context"
	"errors"
	"os"

	"github.com/fadilmartias/dilz_code/apps/backend/app/models"
	"github.com/fadilmartias/dilz_code/apps/backend/app/repositories"
	"github.com/fadilmartias/dilz_code/apps/backend/app/requests"
	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v3"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
)

type AuthService struct {
	DB                 *gorm.DB
	UserRepository     *repositories.UserRepository
	ActivityLogService *ActivityLogService
	Redis              *redis.Client
	TokenIssuer        *AuthTokenIssuer
	OTPService         *AuthOTPService
	OAuthService       *AuthOAuthService
}

func NewAuthService(db *gorm.DB, redis *redis.Client, userRepository *repositories.UserRepository, activityLogService *ActivityLogService) *AuthService {
	tokenIssuer := NewAuthTokenIssuer()
	otpService := NewAuthOTPService(redis, userRepository)
	oAuthService := NewAuthOAuthService(db, userRepository, activityLogService, tokenIssuer)

	return &AuthService{
		DB:                 db,
		UserRepository:     userRepository,
		ActivityLogService: activityLogService,
		Redis:              redis,
		TokenIssuer:        tokenIssuer,
		OTPService:         otpService,
		OAuthService:       oAuthService,
	}
}

var Err2FARequired = errors.New("2FA required")

func (s *AuthService) VerifyRecaptcha(token string, action string) error {
	data := map[string]any{
		"event": map[string]any{
			"token":          token,
			"expectedAction": action,
			"siteKey":        os.Getenv("GOOGLE_RECAPTCHA_SITE_KEY"),
		},
	}

	resp, err := utils.Http().WithJSON(data).Post("https://recaptchaenterprise.googleapis.com/v1/projects/nifty-jet-464704-i6/assessments?key=" + os.Getenv("GOOGLE_RECAPTCHA_API_KEY"))
	if err != nil {
		return err
	}

	if resp.StatusCode() != fiber.StatusOK {
		return errors.New("failed to validate recaptcha")
	}
	if gjson.Get(resp.String(), "tokenProperties.action").String() != action {
		return errors.New("recaptcha validation failed")
	}
	if !gjson.Get(resp.String(), "tokenProperties.valid").Bool() {
		return errors.New("recaptcha validation failed")
	}
	if gjson.Get(resp.String(), "riskAnalysis.score").Float() < 0.5 {
		return errors.New("recaptcha validation failed")
	}
	return nil
}

func (s *AuthService) CreateNewUser(trx *gorm.DB, input requests.RegisterInput, tenantID string) (*models.User, error) {
	return s.OAuthService.CreateNewUser(trx, input, tenantID)
}

func (s *AuthService) Register(c fiber.Ctx, input requests.RegisterInput) error {
	trx := s.DB.Begin()
	if trx.Error != nil {
		return trx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			trx.Rollback()
		}
	}()

	if _, err := s.CreateNewUser(trx, input, "1"); err != nil {
		trx.Rollback()
		return err
	}

	return trx.Commit().Error
}

func (s *AuthService) Login(ctx context.Context, credential, password string) (*models.User, string, string, string, error) {
	user, err := s.UserRepository.FindByEmail(credential)
	if err != nil {
		return nil, "", "", "", err
	}
	if err := user.CheckPassword(password); err != nil {
		return nil, "", "", "", err
	}

	if user.TOTPSecret != nil {
		tempToken, err := s.TokenIssuer.IssueTempToken(user.Email)
		if err != nil {
			return nil, "", "", "", err
		}
		return user, "", "", tempToken, Err2FARequired
	}

	accessToken, refreshToken, err := s.TokenIssuer.IssueTokenPair(user, false)
	if err != nil {
		return nil, "", "", "", err
	}

	user.RefreshToken = &refreshToken
	if _, err := s.UserRepository.UpdateByModel(user); err != nil {
		return nil, "", "", "", err
	}

	entry := models.ActivityLog{
		LogName:     "auth",
		Desc:        "User login",
		Event:       "login",
		SubjectType: "user",
		SubjectID:   &user.ID,
		CauserType:  "user",
		CauserID:    &user.ID,
		Properties:  nil,
	}
	_ = s.ActivityLogService.Log(ctx, &entry)

	return user, accessToken, refreshToken, "", nil
}

func (s *AuthService) Login2FA(ctx context.Context, tempToken, code string, isLogin bool) (*models.User, string, string, error) {
	claimsTemp, err := utils.ValidateToken(tempToken)
	if err != nil {
		return nil, "", "", err
	}

	email := claimsTemp["email"].(string)
	user, err := s.UserRepository.FindByEmail(email)
	if err != nil {
		return nil, "", "", err
	}
	if err := s.VerifyUserTOTPSecret(email, code, isLogin); err != nil {
		return nil, "", "", err
	}

	accessToken, refreshToken, err := s.TokenIssuer.IssueTokenPair(user, true)
	if err != nil {
		return nil, "", "", err
	}

	user.RefreshToken = &refreshToken
	if _, err := s.UserRepository.UpdateByModel(user); err != nil {
		return nil, "", "", err
	}

	entry := models.ActivityLog{
		LogName:     "auth",
		Desc:        "User login with 2FA",
		Event:       "login",
		SubjectType: "user",
		SubjectID:   &user.ID,
		CauserType:  "user",
		CauserID:    &user.ID,
		Properties:  nil,
	}
	_ = s.ActivityLogService.Log(ctx, &entry)

	return user, accessToken, refreshToken, nil
}

func (s *AuthService) Register2FA(email string) (string, string, error) {
	return s.OTPService.Register2FA(email)
}

func (s *AuthService) VerifyUserTOTPSecret(email string, code string, isLogin bool) error {
	return s.OTPService.VerifyUserTOTPSecret(email, code, isLogin)
}

func (s *AuthService) Disable2FA(id string) error {
	return s.OTPService.Disable2FA(id)
}

func (s *AuthService) SendOTP(ctx context.Context, targetType string, target string, subject string, name string) error {
	return s.OTPService.SendOTP(ctx, targetType, target, subject, name)
}

func (s *AuthService) VerifyOTP(ctx context.Context, target string, subject string, otp string) error {
	return s.OTPService.VerifyOTP(ctx, target, subject, otp)
}

func (s *AuthService) HandleOAuth(ctx context.Context, provider, id, email, name string) (*models.User, string, string, string, bool, error) {
	return s.OAuthService.HandleOAuth(ctx, provider, id, email, name)
}
