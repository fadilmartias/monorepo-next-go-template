package controllers_v1

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/fadilmartias/dilz_code/apps/backend/app/client"
	"github.com/fadilmartias/dilz_code/apps/backend/app/models"
	"github.com/fadilmartias/dilz_code/apps/backend/app/requests"
	"github.com/fadilmartias/dilz_code/apps/backend/app/responses"
	"github.com/fadilmartias/dilz_code/apps/backend/app/services"
	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"
	"github.com/fadilmartias/dilz_code/apps/backend/config"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"github.com/tidwall/gjson"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm" // TAMBAHKAN IMPORT INI
)

type AuthController struct {
	BaseController
	UserService  *services.UserService
	EmailService *services.EmailService
	AuthService  *services.AuthService
	DB           *gorm.DB // Tambahkan ini untuk menyimpan koneksi DB
	Redis        *redis.Client
}

// Ubah fungsi NewAuthController untuk menerima koneksi DB
func NewAuthController(db *gorm.DB, redis *redis.Client, userService *services.UserService, emailService *services.EmailService, authService *services.AuthService) *AuthController {
	return &AuthController{DB: db, Redis: redis, UserService: userService, EmailService: emailService, AuthService: authService}
}

type ValidateRecaptchaInput struct {
	Token  string `json:"token"`
	Action string `json:"action"`
}

func (ctrl *AuthController) VerifyRecaptcha(c fiber.Ctx) error {
	var input ValidateRecaptchaInput
	if err := c.Bind().Body(&input); err != nil {
		return err
	}
	if input.Token == "" || input.Action == "" {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusBadRequest,
			Message: "Token dan action harus diisi",
		})
	}

	if err := ctrl.AuthService.VerifyRecaptcha(input.Token, input.Action); err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusUnauthorized,
			Message: "Validasi recaptcha gagal",
		}, err)
	}

	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    fiber.StatusOK,
		Message: "Validasi recaptcha berhasil",
	})

}

func (ctrl *AuthController) Me(c fiber.Ctx) error {
	user, ok := c.Value("user").(jwt.MapClaims)
	if !ok {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusUnauthorized,
			Message: "Invalid token",
		})
	}
	userId := user["id"].(string)
	userDB, err := ctrl.UserService.FindByID(userId)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusNotFound,
			Message: "User tidak ditemukan",
		}, err)
	}

	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    fiber.StatusOK,
		Message: "User retrieved successfully",
		Data:    userDB,
	})
}

func (ctrl *AuthController) Logout(c fiber.Ctx) error {
	user, ok := c.Value("user").(jwt.MapClaims)
	if !ok {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusUnauthorized,
			Message: "Invalid token",
		})
	}

	var userDB models.User
	if err := ctrl.DB.Where("email = ?", user["email"]).First(&userDB).Error; err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusNotFound,
			Message: "User not found",
		}, err)
	}

	userDB.RefreshToken = nil

	ctrl.clearCookie(c)

	if err := ctrl.DB.Save(&userDB).Error; err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to logout user",
		}, err)
	}

	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    fiber.StatusOK,
		Message: "Logged out successfully",
	})
}

// Register membuat user baru
func (ctrl *AuthController) Register(c fiber.Ctx) error {
	input, err := utils.GetValidatedBody[requests.RegisterInput](c)
	if err != nil {
		return err
	}

	if err := ctrl.AuthService.Register(c, input); err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		}, err)
	}

	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    fiber.StatusCreated,
		Message: "User created successfully",
	})
}

// Login placeholder
func (ctrl *AuthController) Login(c fiber.Ctx) error {

	ctx := c.Context()

	// Parse input dari body
	input, err := utils.GetValidatedBody[requests.LoginInput](c)
	if err != nil {
		return err
	}

	user, accessToken, refreshToken, tempToken, err := ctrl.AuthService.Login(ctx, input.Credential, input.Password)

	// Cek kalau butuh 2FA
	if errors.Is(err, services.Err2FARequired) {
		return utils.SuccessResponse(c, utils.SuccessResponseFormat{
			Code:    fiber.StatusOK,
			Message: "2FA diperlukan",
			Data: fiber.Map{
				"is_2fa_enabled": true,
				"temp_token":     tempToken,
			},
		})
	}

	// Cek error lain
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusUnauthorized,
			Message: "Email atau password salah / user tidak ditemukan",
		}, err)
	}

	ctrl.setCookie(c, accessToken, refreshToken)

	response := responses.UserResponse{
		ID:              &user.ID,
		Name:            &user.Name,
		Phone:           &user.Phone,
		Email:           &user.Email,
		Role:            &user.Role,
		EmailVerifiedAt: user.EmailVerifiedAt,
		TenantID:        &user.TenantID,
		CreatedAt:       &user.CreatedAt,
		UpdatedAt:       &user.UpdatedAt,
	}

	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    fiber.StatusOK,
		Message: "User logged in successfully",
		Data:    response,
	})
}

func (ctrl *AuthController) Login2FA(c fiber.Ctx) error {
	input, err := utils.GetValidatedBody[requests.Login2FARequest](c)
	if err != nil {
		return err
	}
	ctx := c.Context()
	user, accessToken, refreshToken, err := ctrl.AuthService.Login2FA(ctx, input.TempToken, input.Code, true)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusUnauthorized,
			Message: err.Error(),
			Details: fiber.Map{
				"error_code": utils.TranslateErrorCode(err),
			},
		}, err)
	}
	ctrl.setCookie(c, accessToken, refreshToken)

	response := responses.UserResponse{
		ID:              &user.ID,
		Name:            &user.Name,
		Phone:           &user.Phone,
		Email:           &user.Email,
		Role:            &user.Role,
		EmailVerifiedAt: user.EmailVerifiedAt,
		TenantID:        &user.TenantID,
		CreatedAt:       &user.CreatedAt,
		UpdatedAt:       &user.UpdatedAt,
	}
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    fiber.StatusOK,
		Message: "User logged in successfully",
		Data:    response,
	})
}

func (ctrl *AuthController) ForgotPassword(c fiber.Ctx) error {
	input, err := utils.GetValidatedBody[requests.ForgotPasswordInput](c)
	if err != nil {
		return err
	}

	var user = models.User{}
	var passwordResetToken models.PasswordResetToken
	if err := ctrl.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusNotFound,
			Message: "Pengguna tidak ditemukan",
		}, err)
	}

	if err := ctrl.DB.Where("email = ?", input.Email).Where("expired_at > ?", time.Now()).First(&passwordResetToken).Error; err == nil {
		return utils.SuccessResponse(c, utils.SuccessResponseFormat{
			Code:    fiber.StatusOK,
			Message: "Token reset password sudah dikirim",
		})
	}

	// Generate random token
	token := utils.GenerateRandomToken(32)
	tokenHash := utils.HashToken(token)

	passwordResetToken.Email = user.Email
	passwordResetToken.Token = tokenHash
	passwordResetToken.ExpiredAt = time.Now().Add(time.Minute * 5)
	if err := ctrl.DB.Save(&passwordResetToken).Error; err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusInternalServerError,
			Message: "Gagal menyimpan token reset password",
		}, err)
	}

	// Send email
	if err := ctrl.EmailService.SendEmailResetPassword(user.Email, token, user.Name); err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusInternalServerError,
			Message: "Gagal mengirim email reset password",
		}, err)
	}
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    fiber.StatusOK,
		Message: "Token reset password berhasil dikirim",
	})
}

func (ctrl *AuthController) ResetPassword(c fiber.Ctx) error {
	input, err := utils.GetValidatedBody[requests.ResetPasswordInput](c)
	if err != nil {
		return err
	}
	tokenHash := utils.HashToken(input.Token)
	var passwordResetToken models.PasswordResetToken
	if err := ctrl.DB.Where("token = ?", tokenHash).First(&passwordResetToken).Error; err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusNotFound,
			Message: "Password reset token not found",
		}, err)
	}

	if passwordResetToken.UsedAt != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusConflict,
			Message: "Password reset token already used",
			Details: nil,
		})
	}

	user := models.User{}
	if err := ctrl.DB.Where("email = ?", passwordResetToken.Email).First(&user).Error; err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusNotFound,
			Message: "Pengguna tidak ditemukan",
		}, err)
	}

	user.HashPassword(input.Password)
	if err := ctrl.DB.Save(&user).Error; err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusInternalServerError,
			Message: "Gagal menyimpan user",
		}, err)
	}

	now := time.Now()

	passwordResetToken.UsedAt = &now
	ctrl.DB.Save(&passwordResetToken)

	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    fiber.StatusOK,
		Message: "Password reset berhasil",
	})
}

func (ctrl *AuthController) SendEmailVerification(c fiber.Ctx) error {
	id := c.Value("user").(jwt.MapClaims)["id"].(string)
	var user models.User
	if err := ctrl.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusNotFound,
			Message: "Pengguna tidak ditemukan",
			Details: nil,
		})
	}
	if user.EmailVerifiedAt != nil {
		return utils.SuccessResponse(c, utils.SuccessResponseFormat{
			Code:    fiber.StatusOK,
			Message: "Email sudah terverifikasi",
		})
	}
	_, err := ctrl.Redis.Get(c.Context(), client.Key(fmt.Sprintf("email_verification_token:%s", user.Email))).Result()
	if err == nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusConflict,
			Message: "Email verifikasi sudah dikirim",
		})
	}

	jwtToken, _ := utils.GenerateToken(map[string]any{
		"email": user.Email,
	}, time.Minute*5)

	if err := ctrl.Redis.Set(c.Context(), client.Key(fmt.Sprintf("email_verification_token:%s", user.Email)), jwtToken, time.Minute*5).Err(); err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusInternalServerError,
			Message: "Gagal menyimpan token verifikasi email",
		}, err)
	}

	if err := ctrl.EmailService.SendEmailVerification(user.Email, jwtToken, user.Name); err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusInternalServerError,
			Message: "Gagal mengirim email verifikasi",
		}, err)
	}

	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    fiber.StatusOK,
		Message: "Email verifikasi berhasil dikirim",
	})
}

func (ctrl *AuthController) VerifyEmail(c fiber.Ctx) error {
	input, err := utils.GetValidatedBody[requests.VerifyEmailInput](c)
	if err != nil {
		return err
	}
	token := input.Token
	if token == "" {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusUnauthorized,
			Message: "Token tidak valid",
		})
	}

	claims, err := utils.ValidateToken(token)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusUnauthorized,
			Message: "Token tidak valid",
		}, err)
	}

	user := models.User{}
	if err := ctrl.DB.Where("email = ?", claims["email"]).First(&user).Error; err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusNotFound,
			Message: "Pengguna tidak ditemukan",
		}, err)
	}
	if user.EmailVerifiedAt != nil {
		return utils.SuccessResponse(c, utils.SuccessResponseFormat{
			Code:    fiber.StatusOK,
			Message: "Email sudah terverifikasi",
		})
	}
	emailVerifiedAt := time.Now()
	user.EmailVerifiedAt = &emailVerifiedAt
	if result := ctrl.DB.Save(&user); result.Error != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusInternalServerError,
			Message: "Gagal menyimpan user",
		}, result.Error)
	}
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    fiber.StatusOK,
		Message: "Email verified successfully",
	})

}

// RefreshToken mengembalikan token yang baru
func (ctrl *AuthController) RefreshAccessToken(c fiber.Ctx) error {
	token := c.Cookies("refresh_token_" + os.Getenv("APP_ENV"))
	if token == "" {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusBadRequest,
			Message: "Token tidak ditemukan",
		})
	}
	claims, err := utils.ValidateToken(token)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusBadRequest,
			Message: "Token tidak valid",
		}, err)
	}

	user := models.User{}
	if err := ctrl.DB.Where("email = ?", claims["email"]).First(&user).Error; err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusNotFound,
			Message: "Pengguna tidak ditemukan",
		}, err)
	}

	accessToken, err := utils.GenerateToken(map[string]any{
		"id":             user.ID,
		"name":           user.Name,
		"tenant_id":      user.TenantID,
		"phone":          user.Phone,
		"email":          user.Email,
		"role":           user.Role,
		"is_2fa_enabled": user.TOTPSecret != nil,
	}, time.Hour*1)

	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusInternalServerError,
			Message: "Gagal menghasilkan access token",
		}, err)
	}

	refreshToken, err := utils.GenerateToken(map[string]any{
		"id":             user.ID,
		"name":           user.Name,
		"tenant_id":      user.TenantID,
		"phone":          user.Phone,
		"email":          user.Email,
		"role":           user.Role,
		"is_2fa_enabled": user.TOTPSecret != nil,
	}, time.Hour*24)

	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusInternalServerError,
			Message: "Gagal menghasilkan refresh token",
		}, err)
	}
	ctrl.setCookie(c, accessToken, refreshToken)

	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    fiber.StatusOK,
		Message: "Token refreshed successfully",
		Data: fiber.Map{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	})
}

func (ctrl *AuthController) GoogleRedirect(c fiber.Ctx) error {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	redirectURI := os.Getenv("APP_URL") + "/v1/auth/google/callback"
	authURL := fmt.Sprintf(
		"https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=email profile",
		clientID, redirectURI,
	)
	return c.Redirect().To(authURL)
}

func (ctrl *AuthController) GoogleCallback(c fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusBadRequest,
			Message: "Code not found",
		})
	}

	googleConfig := config.LoadGoogleConfig()

	clientID := googleConfig.ClientID
	clientSecret := googleConfig.ClientSecret
	redirectURI := os.Getenv("APP_URL") + "/v1/auth/google/callback"

	// Tukar code ke token
	resp, err := utils.Http().
		WithFormData(map[string]string{
			"code":          code,
			"client_id":     clientID,
			"client_secret": clientSecret,
			"redirect_uri":  redirectURI,
			"grant_type":    "authorization_code",
		}).
		Post("https://oauth2.googleapis.com/token")
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusBadRequest,
			Message: "Error get token",
		}, err)
	}

	accessTokenGoogle := gjson.GetBytes(resp.Body(), "access_token").String()

	// Ambil profile
	userResp, err := utils.Http().
		WithAuthToken(accessTokenGoogle).
		Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusBadRequest,
			Message: "Error get user profile",
		}, err)
	}

	email := gjson.GetBytes(userResp.Body(), "email").String()
	name := gjson.GetBytes(userResp.Body(), "name").String()
	id := gjson.GetBytes(userResp.Body(), "id").String()

	user, accessToken, refreshToken, tempToken, _, err := ctrl.AuthService.HandleOAuth(c.Context(), "google", id, email, name)

	// if isRegister {
	// 	return c.Redirect(os.Getenv("FE_URL") + "/onboarding/" + user.ID)
	// }

	// Cek kalau butuh 2FA
	if errors.Is(err, services.Err2FARequired) {
		return c.Redirect().To(os.Getenv("FE_URL") + "/login/2fa/" + tempToken)
	}

	// Cek error lain
	if err != nil {
		return c.Redirect().To(os.Getenv("FE_URL") + "/login?error=" + url.QueryEscape(err.Error()))
	}

	ctrl.setCookie(c, accessToken, refreshToken)
	if user.Role == "admin" {
		return c.Redirect().To(os.Getenv("FE_URL") + "/admin/dashboard")
	}
	// Redirect ke FE
	return c.Redirect().To(os.Getenv("FE_URL"))
}

func (ctrl *AuthController) DiscordRedirect(c fiber.Ctx) error {
	discordConfig := config.LoadDiscordConfig()
	clientID := discordConfig.ClientID
	redirectURI := os.Getenv("APP_URL") + "/v1/auth/discord/callback"
	authURL := fmt.Sprintf(
		"https://discord.com/api/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=identify+email",
		clientID, redirectURI,
	)
	return c.Redirect().To(authURL)
}

// Callback endpoint
func (ctrl *AuthController) DiscordCallback(c fiber.Ctx) error {
	discordConfig := config.LoadDiscordConfig()
	appConfig := config.LoadAppConfig()
	code := c.Query("code")
	if code == "" {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusBadRequest,
			Message: "missing code",
		})
	}

	redirectURI := appConfig.BaseURL + "/v1/auth/discord/callback"

	// Tukar code ke token
	resp, err := utils.Http().
		WithFormData(map[string]string{
			"code":          code,
			"client_id":     discordConfig.ClientID,
			"client_secret": discordConfig.ClientSecret,
			"redirect_uri":  redirectURI,
			"grant_type":    "authorization_code",
		}).
		Post("https://discord.com/api/oauth2/token")
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusBadRequest,
			Message: "Error get token",
		}, err)
	}

	accessTokenDiscord := gjson.GetBytes(resp.Body(), "access_token").String()

	// Ambil profile
	userResp, err := utils.Http().
		WithAuthToken(accessTokenDiscord).
		Get("https://discord.com/api/users/@me")
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusBadRequest,
			Message: "Error get user profile",
		}, err)
	}

	discordID := gjson.GetBytes(userResp.Body(), "id").String()
	username := gjson.GetBytes(userResp.Body(), "username").String()
	// discriminator := gjson.GetBytes(userResp.Body(), "discriminator").String()
	email := gjson.GetBytes(userResp.Body(), "email").String()
	// avatar := gjson.GetBytes(userResp.Body(), "avatar").String()

	user, accessToken, refreshToken, tempToken, _, err := ctrl.AuthService.HandleOAuth(c.Context(), "discord", discordID, email, username)

	// if isRegister {
	// 	return c.Redirect(os.Getenv("FE_URL") + "/onboarding/" + user.ID)
	// }

	// Cek kalau butuh 2FA
	if errors.Is(err, services.Err2FARequired) {
		return c.Redirect().To(os.Getenv("FE_URL") + "/login/2fa/" + tempToken)
	}

	// Cek error lain
	if err != nil {
		return c.Redirect().To(os.Getenv("FE_URL") + "/login?error=" + url.QueryEscape(err.Error()))
	}
	ctrl.setCookie(c, accessToken, refreshToken)
	if user.Role == "admin" {
		return c.Redirect().To(os.Getenv("FE_URL") + "/admin/dashboard")
	}
	// Redirect ke FE
	return c.Redirect().To(os.Getenv("FE_URL"))

}

func (ctrl *AuthController) GoogleOneTap(c fiber.Ctx) error {
	var body struct {
		Credential string `json:"credential"`
	}
	if err := c.Bind().Body(&body); err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request",
		}, err)
	}

	clientID := os.Getenv("GOOGLE_CLIENT_ID")

	resp, err := utils.Http().
		WithQuery(map[string]string{"id_token": body.Credential}).
		Get("https://oauth2.googleapis.com/tokeninfo")

	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusInternalServerError,
			Message: "Google verify error: " + err.Error(),
		}, err)
	}
	if resp.StatusCode() != 200 {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusUnauthorized,
			Message: "Invalid ID Token",
		}, err)
	}

	// Ambil data dari response
	email := gjson.GetBytes(resp.Body(), "email").String()
	name := gjson.GetBytes(resp.Body(), "name").String()
	aud := gjson.GetBytes(resp.Body(), "aud").String()
	id := gjson.GetBytes(resp.Body(), "id").String()

	// Pastikan audience sesuai CLIENT_ID kamu
	if aud != clientID {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusInternalServerError,
			Message: "Error credential",
		}, err)
	}

	user, accessToken, refreshToken, tempToken, _, err := ctrl.AuthService.HandleOAuth(c.Context(), "google", id, email, name)

	// if isRegister {
	// 	return c.Redirect(os.Getenv("FE_URL") + "/onboarding/" + user.ID)
	// }

	// Cek kalau butuh 2FA
	if errors.Is(err, services.Err2FARequired) {
		return utils.SuccessResponse(c, utils.SuccessResponseFormat{
			Code:    fiber.StatusOK,
			Message: "2FA diperlukan",
			Data: fiber.Map{
				"is_2fa_enabled": true,
				"temp_token":     tempToken,
			},
		})
	}

	// Cek error lain
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusInternalServerError,
			Message: "Error handle user google",
		}, err)
	}
	ctrl.setCookie(c, accessToken, refreshToken)

	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Success",
		Data:    user,
	})
}

func (ctrl *AuthController) FacebookRedirect(c fiber.Ctx) error {
	facebookConfig := config.LoadFacebookConfig()
	appConfig := config.LoadAppConfig()
	redirectURI := appConfig.BaseURL + "/v1/auth/facebook/callback"
	authURL := fmt.Sprintf(
		"https://www.facebook.com/v23.0/dialog/oauth?client_id=%s&redirect_uri=%s&scope=email,public_profile&response_type=code&state=%s",
		facebookConfig.ClientID, url.QueryEscape(redirectURI), "randomstate",
	)
	return c.Redirect().To(authURL)
}

func (ctrl *AuthController) FacebookCallback(c fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusBadRequest,
			Message: "Missing code",
		}, nil)
	}

	facebookConfig := config.LoadFacebookConfig()
	appConfig := config.LoadAppConfig()
	clientID := facebookConfig.ClientID
	clientSecret := facebookConfig.ClientSecret
	redirectURI := appConfig.BaseURL + "/v1/auth/facebook/callback"

	// 1. Tukar code -> access token
	tokenURL := fmt.Sprintf(
		"https://graph.facebook.com/v23.0/oauth/access_token?client_id=%s&redirect_uri=%s&client_secret=%s&code=%s",
		clientID, url.QueryEscape(redirectURI), clientSecret, code,
	)

	resp, err := utils.Http().Get(tokenURL)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusInternalServerError,
			Message: "Facebook token request error: " + err.Error(),
		}, err)
	}
	if resp.StatusCode() != 200 {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusUnauthorized,
			Message: "Invalid Facebook token response",
		}, nil)
	}

	accessToken := gjson.GetBytes(resp.Body(), "access_token").String()
	if accessToken == "" {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusUnauthorized,
			Message: "Missing access_token in response",
		}, nil)
	}

	// 2. Ambil user info dari Facebook Graph API
	userInfoURL := fmt.Sprintf("https://graph.facebook.com/me?fields=id,name,email&access_token=%s", accessToken)
	userResp, err := utils.Http().Get(userInfoURL)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusInternalServerError,
			Message: "Facebook userinfo request error: " + err.Error(),
		}, err)
	}
	if userResp.StatusCode() != 200 {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusUnauthorized,
			Message: "Invalid Facebook userinfo response",
		}, nil)
	}

	id := gjson.GetBytes(userResp.Body(), "id").String()
	name := gjson.GetBytes(userResp.Body(), "name").String()
	email := gjson.GetBytes(userResp.Body(), "email").String()

	if id == "" {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusUnauthorized,
			Message: "Failed to get Facebook user id",
		}, nil)
	}

	// 3. Simpan / update user di DB
	user, accessToken, refreshToken, tempToken, _, err := ctrl.AuthService.HandleOAuth(c.Context(), "facebook", id, email, name)

	// if isRegister {
	// 	return c.Redirect(os.Getenv("FE_URL") + "/onboarding/" + user.ID)
	// }

	// Cek kalau butuh 2FA
	if errors.Is(err, services.Err2FARequired) {
		return c.Redirect().To(os.Getenv("FE_URL") + "/login/2fa/" + tempToken)
	}

	// Cek error lain
	if err != nil {
		return c.Redirect().To(os.Getenv("FE_URL") + "/login?error=" + url.QueryEscape(err.Error()))
	}
	ctrl.setCookie(c, accessToken, refreshToken)
	if user.Role == "admin" {
		return c.Redirect().To(os.Getenv("FE_URL") + "/admin/dashboard")
	}

	return c.Redirect().To(os.Getenv("FE_URL"))
}

func (ctrl *AuthController) SteamRedirect(c fiber.Ctx) error {
	appConfig := config.LoadAppConfig()
	redirectURI := appConfig.BaseURL + "/v1/auth/steam/callback"
	authURL := fmt.Sprintf(
		"https://steamcommunity.com/openid/login?openid.return_to=%s&openid.realm=%s",
		url.QueryEscape(redirectURI), url.QueryEscape(appConfig.BaseURL),
	)
	return c.Redirect().To(authURL)
}

func (ctrl *AuthController) SteamCallback(c fiber.Ctx) error {
	// 1. Kirim balik semua param openid ke Steam untuk verifikasi
	form := url.Values{}
	for k, v := range c.Queries() {
		if strings.HasPrefix(k, "openid.") {
			form.Set(k, v)
		}
	}
	form.Set("openid.mode", "check_authentication")
	formData := make(map[string]string, len(form))
	for key, values := range form {
		if len(values) > 0 {
			formData[key] = values[0]
		}
	}

	resp, err := utils.Http().
		WithFormData(formData).
		Post("https://steamcommunity.com/openid/login")
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusInternalServerError,
			Message: "Steam verification failed: " + err.Error(),
		}, err)
	}
	if resp.StatusCode() != 200 || !strings.Contains(string(resp.Body()), "is_valid:true") {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusUnauthorized,
			Message: "Invalid Steam login",
		}, nil)
	}

	// 2. Ambil SteamID dari openid.claimed_id
	claimedID := c.FormValue("openid.claimed_id")
	if claimedID == "" {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusBadRequest,
			Message: "Missing claimed_id",
		}, nil)
	}
	steamID := claimedID[strings.LastIndex(claimedID, "/")+1:]

	// 3. Ambil profil user dari Steam Web API (opsional)
	apiKey := os.Getenv("STEAM_API_KEY")
	var name string
	if apiKey != "" {
		profileURL := fmt.Sprintf(
			"https://api.steampowered.com/ISteamUser/GetPlayerSummaries/v2/?key=%s&steamids=%s",
			apiKey, steamID,
		)
		profileResp, err := utils.Http().Get(profileURL)
		if err == nil && profileResp.StatusCode() == 200 {
			name = gjson.GetBytes(profileResp.Body(), "response.players.0.personaname").String()
		}
	}
	if name == "" {
		name = "SteamUser_" + steamID // fallback kalau API kosong
	}

	// 4. Simpan / update user via handleUserOAuth
	user, accessToken, refreshToken, tempToken, _, err := ctrl.AuthService.HandleOAuth(c.Context(), "steam", steamID, "", name)

	// if isRegister {
	// 	return c.Redirect(os.Getenv("FE_URL") + "/onboarding/" + user.ID)
	// }

	// Cek kalau butuh 2FA
	if errors.Is(err, services.Err2FARequired) {
		return c.Redirect().To(os.Getenv("FE_URL") + "/login/2fa/" + tempToken)
	}

	// Cek error lain
	if err != nil {
		return c.Redirect().To(os.Getenv("FE_URL") + "/login?error=" + url.QueryEscape(err.Error()))
	}
	ctrl.setCookie(c, accessToken, refreshToken)
	if user.Role == "admin" {
		return c.Redirect().To(os.Getenv("FE_URL") + "/admin/dashboard")
	}
	return c.Redirect().To(os.Getenv("FE_URL"))
}

func (ctrl *AuthController) TwitchRedirect(c fiber.Ctx) error {
	twitchConfig := config.LoadTwitchConfig()
	appConfig := config.LoadAppConfig()

	clientID := twitchConfig.ClientID
	redirectURI := appConfig.BaseURL + "/v1/auth/twitch/callback"

	authURL := fmt.Sprintf(
		"https://id.twitch.tv/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=user:read:email",
		clientID, url.QueryEscape(redirectURI),
	)

	return c.Redirect().To(authURL)
}

func (ctrl *AuthController) TwitchCallback(c fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return c.Status(400).SendString("missing code")
	}

	twitchConfig := config.LoadTwitchConfig()
	appConfig := config.LoadAppConfig()

	clientID := twitchConfig.ClientID
	clientSecret := twitchConfig.ClientSecret
	redirectURI := appConfig.BaseURL + "/v1/auth/twitch/callback"

	// 1. Tukar code -> access token
	tokenResp, err := utils.Http().
		WithHeader("Content-Type", "application/x-www-form-urlencoded").
		WithFormData(map[string]string{
			"client_id":     clientID,
			"client_secret": clientSecret,
			"code":          code,
			"grant_type":    "authorization_code",
			"redirect_uri":  redirectURI,
		}).
		Post("https://id.twitch.tv/oauth2/token")

	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusInternalServerError,
			Message: "Twitch token exchange error: " + err.Error(),
		}, err)
	}
	if tokenResp.StatusCode() != 200 {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusUnauthorized,
			Message: "Invalid Twitch authorization code",
		}, nil)
	}

	accessToken := gjson.GetBytes(tokenResp.Body(), "access_token").String()
	if accessToken == "" {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusUnauthorized,
			Message: "Access token missing",
		}, nil)
	}

	// 2. Ambil data user
	userResp, err := utils.Http().
		WithHeaders(map[string]string{
			"Authorization": "Bearer " + accessToken,
			"Client-Id":     clientID,
		}).
		Get("https://api.twitch.tv/helix/users")

	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusInternalServerError,
			Message: "Twitch user fetch error: " + err.Error(),
		}, err)
	}
	if userResp.StatusCode() != 200 {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusUnauthorized,
			Message: "Failed to fetch Twitch user",
		}, nil)
	}

	// Twitch response → "data":[{...}]
	id := gjson.GetBytes(userResp.Body(), "data.0.id").String()
	email := gjson.GetBytes(userResp.Body(), "data.0.email").String()
	name := gjson.GetBytes(userResp.Body(), "data.0.display_name").String()

	if id == "" {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusUnauthorized,
			Message: "Missing Twitch user ID",
		}, nil)
	}

	// 3. Integrasi ke sistem
	user, accessToken, refreshToken, tempToken, _, err := ctrl.AuthService.HandleOAuth(c.Context(), "twitch", id, email, name)

	// if isRegister {
	// 	return c.Redirect(os.Getenv("FE_URL") + "/onboarding/" + user.ID)
	// }

	// Cek kalau butuh 2FA
	if errors.Is(err, services.Err2FARequired) {
		return c.Redirect().To(os.Getenv("FE_URL") + "/login/2fa/" + tempToken)
	}

	// Cek error lain
	if err != nil {
		return c.Redirect().To(os.Getenv("FE_URL") + "/login?error=" + url.QueryEscape(err.Error()))
	}
	ctrl.setCookie(c, accessToken, refreshToken)
	if user.Role == "admin" {
		return c.Redirect().To(os.Getenv("FE_URL") + "/admin/dashboard")
	}
	return c.Redirect().To(os.Getenv("FE_URL"))
}

func (ctrl *AuthController) Register2FA(c fiber.Ctx) error {

	user := c.Value("user").(jwt.MapClaims)

	// Validasi email
	if user["email"] == "" {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusBadRequest,
			Message: "Email tidak boleh kosong",
		}, nil)
	}

	secret, qrString, err := ctrl.AuthService.Register2FA(user["email"].(string))
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusInternalServerError,
			Message: "Gagal menyimpan kunci rahasia",
		}, err)
	}

	// Kembalikan secret key dan URL untuk QR code
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Success",
		Data: fiber.Map{
			"totp_secret": secret,
			"qr_string":   qrString,
		},
	})
}

func (ctrl *AuthController) Verify2FA(c fiber.Ctx) error {
	user := c.Value("user").(jwt.MapClaims)

	input, err := utils.GetValidatedBody[requests.Verify2FARequest](c)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusBadRequest,
			Message: "Gagal parsing request body",
		}, err)
	}

	// Validasi input
	if user["email"] == "" || input.Code == "" {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusBadRequest,
			Message: "Email dan kode harus diisi",
		}, nil)
	}

	// Ambil kunci rahasia dari database
	if err := ctrl.AuthService.VerifyUserTOTPSecret(user["email"].(string), input.Code, false); err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusNotFound,
			Message: "User tidak ditemukan atau token tidak valid",
		}, err)
	}

	// Kode valid, kembalikan respons sukses
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Sukses verifikasi 2FA",
	})
}

func (ctrl *AuthController) Disable2FA(c fiber.Ctx) error {
	user := c.Value("user").(jwt.MapClaims)
	if err := ctrl.AuthService.Disable2FA(user["id"].(string)); err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusInternalServerError,
			Message: "Gagal menghapus 2FA",
		}, err)
	}
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Sukses menghapus 2FA",
	})
}

func (ctrl *AuthController) SendOTP(c fiber.Ctx) error {
	user := c.Value("user")
	input, err := utils.GetValidatedBody[requests.SendOTPRequest](c)
	if err != nil {
		return err
	}
	var target string

	if user != nil {
		switch input.TargetType {
		case "email":
			target = user.(jwt.MapClaims)["email"].(string)
		case "wa":
			target = user.(jwt.MapClaims)["phone"].(string)
		}
	} else {
		target = input.Target
	}

	userName := target
	if user != nil {
		userName = user.(jwt.MapClaims)["name"].(string)
	}

	if target == "" {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusBadRequest,
			Message: "Email atau nomor telepon harus diisi",
		}, nil)
	}

	if err := ctrl.AuthService.SendOTP(c.Context(), input.TargetType, target, input.Subject, userName); err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusInternalServerError,
			Message: "Gagal mengirim OTP",
		}, err)
	}
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Success",
	})
}

func (ctrl *AuthController) VerifyOTP(c fiber.Ctx) error {
	user := c.Value("user")
	input, err := utils.GetValidatedBody[requests.VerifyOTPRequest](c)
	if err != nil {
		return err
	}
	var target string

	if user != nil {
		switch input.TargetType {
		case "email":
			target = user.(jwt.MapClaims)["email"].(string)
		case "wa":
			target = user.(jwt.MapClaims)["phone"].(string)
		}
	} else {
		target = input.Target
	}

	if target == "" {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusBadRequest,
			Message: "Email atau nomor telepon harus diisi",
		}, nil)
	}

	if err := ctrl.AuthService.VerifyOTP(c.Context(), target, input.Subject, input.Otp); err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusNotFound,
			Message: "Kode OTP tidak valid",
		}, err)
	}
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Success",
	})
}

func (ctrl *AuthController) setCookie(c fiber.Ctx, accessToken string, refreshToken string) {
	accessCookie := fiber.Cookie{
		Name:     "access_token_" + os.Getenv("APP_ENV"),
		Value:    accessToken,
		Expires:  time.Now().Add(time.Hour * 1),
		HTTPOnly: true,
		Domain:   os.Getenv("FE_DOMAIN"),
		SameSite: fiber.CookieSameSiteNoneMode,
		Secure:   true,
	}

	refreshCookie := fiber.Cookie{
		Name:     "refresh_token_" + os.Getenv("APP_ENV"),
		Value:    refreshToken,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
		Domain:   os.Getenv("FE_DOMAIN"),
		SameSite: fiber.CookieSameSiteNoneMode,
		Secure:   true,
	}
	c.Cookie(&accessCookie)
	c.Cookie(&refreshCookie)
}

func (ctrl *AuthController) clearCookie(c fiber.Ctx) {
	// Hapus access_token
	c.Cookie(&fiber.Cookie{
		Name:     "access_token_" + os.Getenv("APP_ENV"),
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // Expired
		HTTPOnly: true,
		Domain:   os.Getenv("FE_DOMAIN"),
		SameSite: fiber.CookieSameSiteNoneMode,
		Secure:   true,
	})

	// Hapus refresh_token
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token_" + os.Getenv("APP_ENV"),
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // Expired
		HTTPOnly: true,
		Domain:   os.Getenv("FE_DOMAIN"),
		SameSite: fiber.CookieSameSiteNoneMode,
		Secure:   true,
	})
}

// fiber:context-methods migrated
