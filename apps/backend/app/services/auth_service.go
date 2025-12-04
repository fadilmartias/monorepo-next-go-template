package services

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/fadilmartias/dilz_code/apps/backend/app/models"
	"github.com/fadilmartias/dilz_code/apps/backend/app/repositories"
	"github.com/fadilmartias/dilz_code/apps/backend/app/requests"
	"github.com/fadilmartias/dilz_code/apps/backend/app/usecases"
	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"
	"github.com/fadilmartias/dilz_code/apps/backend/config"
	"github.com/gofiber/fiber/v2"
	"github.com/pquerna/otp/totp"
	"github.com/tidwall/gjson"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

type AuthService struct {
	DB                 *gorm.DB
	UserRepository     *repositories.UserRepository
	ActivityLogService *ActivityLogService
	Redis              *config.RedisClient
}

func NewAuthService(db *gorm.DB, redis *config.RedisClient, userRepository *repositories.UserRepository, activityLogService *ActivityLogService) *AuthService {
	return &AuthService{DB: db, UserRepository: userRepository, Redis: redis, ActivityLogService: activityLogService}
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
	const SettingKeyReferralPoint = "referral_point"

	bytes, err := bcrypt.GenerateFromPassword([]byte(input.Password), 14)
	if err != nil {
		return nil, err
	}
	input.Password = string(bytes)

	newUser := models.User{
		Name:     input.Name,
		TenantID: tenantID,
		Email:    input.Email,
		Phone:    input.Phone,
		Password: input.Password,
		Role:     "user",
	}

	if err := s.UserRepository.WithTx(trx).Create(&newUser); err != nil {
		return nil, err
	}

	return &newUser, nil
}

func (s *AuthService) Register(c *fiber.Ctx, input requests.RegisterInput) error {
	trx := s.DB.Begin()
	if trx.Error != nil {
		return trx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			trx.Rollback()
		}
	}()

	_, err := s.CreateNewUser(trx, input, "1")
	if err != nil {
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

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", "", "", err
	}

	is2FAEnabled := user.TOTPSecret != nil
	if is2FAEnabled {
		claims := map[string]any{
			"email": user.Email,
		}
		tempToken, err := utils.GenerateToken(claims, time.Minute*5)
		if err != nil {
			return nil, "", "", "", err
		}
		return user, "", "", tempToken, Err2FARequired
	}

	claims := map[string]any{
		"id":             user.ID,
		"name":           user.Name,
		"tenant_id":      user.TenantID,
		"phone":          user.Phone,
		"email":          user.Email,
		"role":           user.Role,
		"is_2fa_enabled": false,
	}

	accessToken, err := utils.GenerateToken(claims, time.Hour*1)
	if err != nil {
		return nil, "", "", "", err
	}

	refreshToken, err := utils.GenerateToken(claims, time.Hour*24)
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
		// ini totp salah, balikin error biasa
		return nil, "", "", err
	}
	claims := map[string]any{
		"id":             user.ID,
		"name":           user.Name,
		"tenant_id":      user.TenantID,
		"phone":          user.Phone,
		"email":          user.Email,
		"role":           user.Role,
		"is_2fa_enabled": true,
	}
	accessToken, err := utils.GenerateToken(claims, time.Hour*1)
	if err != nil {
		return nil, "", "", err
	}
	refreshToken, err := utils.GenerateToken(claims, time.Hour*24)
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
	// Generate kunci TOTP
	appConfig := config.LoadAppConfig()
	var issuer string
	if appConfig.Env != "production" {
		issuer = appConfig.Name + " " + cases.Title(language.English, cases.NoLower).String(appConfig.Env)
	} else {
		issuer = appConfig.Name
	}
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: email, // Gunakan email sebagai AccountName
	})
	if err != nil {
		return "", "", err
	}
	s.Redis.Set(context.Background(), fmt.Sprintf("totp:%s", email), key.Secret(), 5*time.Minute)
	return key.Secret(), key.URL(), nil
}

func (s *AuthService) VerifyUserTOTPSecret(email string, code string, isLogin bool) error {
	var secret string
	if !isLogin {
		secretRedis, err := s.Redis.Get(context.Background(), fmt.Sprintf("totp:%s", email))
		if err != nil {
			return err
		}
		secret = secretRedis
	} else {
		userDB, err := s.UserRepository.FindByEmail(email)
		secretDB := userDB.TOTPSecret
		if err != nil {
			return err
		}
		if secretDB == nil {
			return errors.New("totp_secret not found")
		}
		secret = *secretDB
	}
	if secret == "" {
		return errors.New("totp_secret not found")
	}
	// Verifikasi kode TOTP
	valid := totp.Validate(code, secret)
	if !valid {
		return errors.New("kode TOTP tidak valid")
	}
	// Hapus kunci TOTP dari Redis
	if err := s.Redis.Del(context.Background(), fmt.Sprintf("totp:%s", email)); err != nil {
		return err
	}
	userDB, err := s.UserRepository.FindByEmail(email)
	if err != nil {
		return err
	}
	if userDB == nil {
		return errors.New("user not found")
	}
	if userDB.TOTPSecret == nil {
		userDB.TOTPSecret = &secret
		s.UserRepository.UpdateWithModel(userDB, map[string]any{"totp_secret": secret})
	}
	return nil
}

func (s *AuthService) Disable2FA(id string) error {
	_, err := s.UserRepository.UpdateByID(id, map[string]any{"totp_secret": nil})
	return err
}

func (s *AuthService) SendOTP(ctx context.Context, targetType string, target string, subject string, name string) error {
	otp := strconv.Itoa(utils.GenerateRandomNumber(6))
	s.Redis.Set(ctx, fmt.Sprintf("otp:%s:%s", subject, target), otp, 5*time.Minute)
	otpUsecase := usecases.NewSendOTPUsecase()
	otpUsecase.Execute(targetType, target, subject, otp, name)
	return nil
}

func (s *AuthService) VerifyOTP(ctx context.Context, target string, subject string, otp string) error {
	storedOtp, err := s.Redis.Get(ctx, fmt.Sprintf("otp:%s:%s", subject, target))
	if err != nil {
		return err
	}
	if storedOtp != otp {
		return errors.New("kode OTP tidak valid")
	}
	return nil
}

func (s *AuthService) HandleOAuth(ctx context.Context, provider, id, email, name string) (*models.User, string, string, string, bool, error) {
	var isRegister bool
	user, err := s.UserRepository.FindByEmail(email)
	if err != nil {
		// --- Jika user belum ada, buat baru ---
		if errors.Is(err, gorm.ErrRecordNotFound) {
			isRegister = true
			trx := s.DB.Begin()
			if trx.Error != nil {
				return nil, "", "", "", isRegister, trx.Error
			}
			defer func() {
				if r := recover(); r != nil {
					trx.Rollback()
				}
			}()

			// Buat input dummy (karena OAuth tidak pakai password & referral code)
			input := requests.RegisterInput{
				Name:     name,
				Email:    email,
				Phone:    "",
				Password: strconv.Itoa(utils.GenerateRandomNumber(16)),
			}

			// Gunakan logic modular dari Register()
			newUser, err := s.CreateNewUser(trx, input, "T1")
			if err != nil {
				trx.Rollback()
				return nil, "", "", "", isRegister, err
			}

			// Tambahkan info provider ID
			switch provider {
			case "google":
				newUser.GoogleID = models.NewNullString(id)
			case "discord":
				newUser.DiscordID = models.NewNullString(id)
			case "facebook":
				newUser.FacebookID = models.NewNullString(id)
			case "steam":
				newUser.SteamID = models.NewNullString(id)
			case "twitch":
				newUser.TwitchID = models.NewNullString(id)
			case "apple":
				newUser.AppleID = models.NewNullString(id)
			}

			emailVerifiedAt := time.Now()
			newUser.EmailVerifiedAt = &emailVerifiedAt

			// Update user dengan data provider dan verifikasi email
			if err := trx.Save(&newUser).Error; err != nil {
				trx.Rollback()
				return nil, "", "", "", isRegister, err
			}

			if err := trx.Commit().Error; err != nil {
				return nil, "", "", "", isRegister, err
			}

			user = newUser
		} else {
			return nil, "", "", "", isRegister, err
		}
	} else {
		// --- Jika user sudah ada, update ID provider ---
		switch provider {
		case "google":
			user.GoogleID = models.NewNullString(id)
		case "discord":
			user.DiscordID = models.NewNullString(id)
		case "facebook":
			user.FacebookID = models.NewNullString(id)
		case "steam":
			user.SteamID = models.NewNullString(id)
		case "twitch":
			user.TwitchID = models.NewNullString(id)
		case "apple":
			user.AppleID = models.NewNullString(id)
		}

		if _, err := s.UserRepository.UpdateByModel(user); err != nil {
			return nil, "", "", "", isRegister, err
		}
	}

	// --- 2FA Check ---
	is2FAEnabled := user.TOTPSecret != nil
	if is2FAEnabled {
		claims := map[string]any{
			"email": user.Email,
		}
		tempToken, err := utils.GenerateToken(claims, time.Minute*5)
		if err != nil {
			return nil, "", "", "", isRegister, err
		}
		return user, "", "", tempToken, isRegister, Err2FARequired
	}

	// --- Generate Access & Refresh Token ---
	claims := map[string]any{
		"id":             user.ID,
		"name":           user.Name,
		"tenant_id":      user.TenantID,
		"phone":          user.Phone,
		"email":          user.Email,
		"role":           user.Role,
		"is_2fa_enabled": user.TOTPSecret != nil,
	}

	accessToken, err := utils.GenerateToken(claims, time.Hour*1)
	if err != nil {
		return nil, "", "", "", isRegister, err
	}

	refreshToken, err := utils.GenerateToken(claims, time.Hour*24)
	if err != nil {
		return nil, "", "", "", isRegister, err
	}

	user.RefreshToken = &refreshToken
	if _, err := s.UserRepository.UpdateByModel(user); err != nil {
		return nil, "", "", "", isRegister, err
	}

	entry := models.ActivityLog{
		LogName:     "auth",
		Desc:        "Login with " + provider,
		Event:       "login",
		SubjectType: "user",
		SubjectID:   &user.ID,
		CauserType:  "user",
		CauserID:    &user.ID,
		Properties:  nil,
	}
	_ = s.ActivityLogService.Log(ctx, &entry)

	return user, accessToken, refreshToken, "", isRegister, nil
}
