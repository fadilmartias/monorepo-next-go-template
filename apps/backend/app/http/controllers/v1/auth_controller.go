package controllers_v1

import (
	"errors"
	"net/url"
	"os"
	"time"

	"github.com/fadilmartias/dilz_code/apps/backend/app/requests"
	"github.com/fadilmartias/dilz_code/apps/backend/app/responses"
	"github.com/fadilmartias/dilz_code/apps/backend/app/services"
	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type AuthController struct {
	BaseController
	AuthService     *services.AuthService
	AuthHTTPService *services.AuthHTTPService
}

func NewAuthController(db *gorm.DB, redisClient any, userService *services.UserService, emailService *services.EmailService, authService *services.AuthService) *AuthController {
	_ = redisClient
	return &AuthController{
		AuthService:     authService,
		AuthHTTPService: services.NewAuthHTTPService(authService, userService, emailService, authService.DB, authService.Redis),
	}
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
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: fiber.StatusBadRequest, Message: "Token dan action harus diisi"})
	}
	if err := ctrl.AuthService.VerifyRecaptcha(input.Token, input.Action); err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: fiber.StatusUnauthorized, Message: "Validasi recaptcha gagal"}, err)
	}
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{Code: fiber.StatusOK, Message: "Validasi recaptcha berhasil"})
}

func (ctrl *AuthController) Me(c fiber.Ctx) error {
	user, ok := c.Value("user").(jwt.MapClaims)
	if !ok {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: fiber.StatusUnauthorized, Message: "Invalid token"})
	}
	userDB, err := ctrl.AuthHTTPService.Me(user)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: fiber.StatusNotFound, Message: "User tidak ditemukan"}, err)
	}
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{Code: fiber.StatusOK, Message: "User retrieved successfully", Data: userDB})
}

func (ctrl *AuthController) Logout(c fiber.Ctx) error {
	user, ok := c.Value("user").(jwt.MapClaims)
	if !ok {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: fiber.StatusUnauthorized, Message: "Invalid token"})
	}
	if err := ctrl.AuthHTTPService.Logout(user); err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: fiber.StatusInternalServerError, Message: "Failed to logout user"}, err)
	}
	ctrl.clearCookie(c)
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{Code: fiber.StatusOK, Message: "Logged out successfully"})
}

func (ctrl *AuthController) Register(c fiber.Ctx) error {
	input, err := utils.GetValidatedBody[requests.RegisterInput](c)
	if err != nil {
		return err
	}
	if err := ctrl.AuthService.Register(c, input); err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: fiber.StatusInternalServerError, Message: err.Error()}, err)
	}
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{Code: fiber.StatusCreated, Message: "User created successfully"})
}

func (ctrl *AuthController) Login(c fiber.Ctx) error {
	ctx := c.Context()
	input, err := utils.GetValidatedBody[requests.LoginInput](c)
	if err != nil {
		return err
	}
	user, accessToken, refreshToken, tempToken, err := ctrl.AuthService.Login(ctx, input.Credential, input.Password)
	if errors.Is(err, services.Err2FARequired) {
		return utils.SuccessResponse(c, utils.SuccessResponseFormat{Code: fiber.StatusOK, Message: "2FA diperlukan", Data: fiber.Map{"is_2fa_enabled": true, "temp_token": tempToken}})
	}
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: fiber.StatusUnauthorized, Message: "Email atau password salah / user tidak ditemukan"}, err)
	}
	ctrl.setCookie(c, accessToken, refreshToken)
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{Code: fiber.StatusOK, Message: "User logged in successfully", Data: responses.UserResponse{ID: &user.ID, Name: &user.Name, Phone: &user.Phone, Email: &user.Email, Role: &user.Role, EmailVerifiedAt: user.EmailVerifiedAt, TenantID: &user.TenantID, CreatedAt: &user.CreatedAt, UpdatedAt: &user.UpdatedAt}})
}

func (ctrl *AuthController) Login2FA(c fiber.Ctx) error {
	input, err := utils.GetValidatedBody[requests.Login2FARequest](c)
	if err != nil {
		return err
	}
	user, accessToken, refreshToken, err := ctrl.AuthService.Login2FA(c.Context(), input.TempToken, input.Code, true)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: fiber.StatusUnauthorized, Message: err.Error(), Details: fiber.Map{"error_code": utils.TranslateErrorCode(err)}}, err)
	}
	ctrl.setCookie(c, accessToken, refreshToken)
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{Code: fiber.StatusOK, Message: "User logged in successfully", Data: responses.UserResponse{ID: &user.ID, Name: &user.Name, Phone: &user.Phone, Email: &user.Email, Role: &user.Role, EmailVerifiedAt: user.EmailVerifiedAt, TenantID: &user.TenantID, CreatedAt: &user.CreatedAt, UpdatedAt: &user.UpdatedAt}})
}

func (ctrl *AuthController) ForgotPassword(c fiber.Ctx) error {
	input, err := utils.GetValidatedBody[requests.ForgotPasswordInput](c)
	if err != nil {
		return err
	}
	if _, err := ctrl.AuthHTTPService.ForgotPassword(input.Email); err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: fiber.StatusInternalServerError, Message: "Gagal memproses forgot password"}, err)
	}
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{Code: fiber.StatusOK, Message: "Token reset password berhasil dikirim"})
}

func (ctrl *AuthController) ResetPassword(c fiber.Ctx) error {
	input, err := utils.GetValidatedBody[requests.ResetPasswordInput](c)
	if err != nil {
		return err
	}
	if err := ctrl.AuthHTTPService.ResetPassword(input.Token, input.Password); err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: fiber.StatusInternalServerError, Message: err.Error()}, err)
	}
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{Code: fiber.StatusOK, Message: "Password reset berhasil"})
}

func (ctrl *AuthController) SendEmailVerification(c fiber.Ctx) error {
	user, ok := c.Value("user").(jwt.MapClaims)
	if !ok {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: fiber.StatusUnauthorized, Message: "Invalid token"})
	}
	if _, err := ctrl.AuthHTTPService.SendEmailVerification(user); err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: fiber.StatusInternalServerError, Message: "Gagal mengirim email verifikasi"}, err)
	}
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{Code: fiber.StatusOK, Message: "Email verifikasi berhasil dikirim"})
}

func (ctrl *AuthController) VerifyEmail(c fiber.Ctx) error {
	input, err := utils.GetValidatedBody[requests.VerifyEmailInput](c)
	if err != nil {
		return err
	}
	user, err := ctrl.AuthHTTPService.VerifyEmail(input.Token)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: fiber.StatusUnauthorized, Message: "Token tidak valid"}, err)
	}
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{Code: fiber.StatusOK, Message: "Email verified successfully", Data: user})
}

func (ctrl *AuthController) RefreshAccessToken(c fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token_" + os.Getenv("APP_ENV"))
	if refreshToken == "" {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: fiber.StatusBadRequest, Message: "Token tidak ditemukan"})
	}
	accessToken, newRefreshToken, user, err := ctrl.AuthHTTPService.RefreshAccessToken(refreshToken)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: fiber.StatusBadRequest, Message: err.Error()}, err)
	}
	ctrl.setCookie(c, accessToken, newRefreshToken)
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{Code: fiber.StatusOK, Message: "Token refreshed successfully", Data: fiber.Map{"access_token": accessToken, "refresh_token": newRefreshToken, "user": user}})
}

func (ctrl *AuthController) GoogleRedirect(c fiber.Ctx) error {
	return c.Redirect().To(ctrl.AuthHTTPService.GoogleRedirectURL())
}

func (ctrl *AuthController) GoogleCallback(c fiber.Ctx) error {
	result, err := ctrl.AuthHTTPService.GoogleCallback(c.Query("code"))
	if errors.Is(err, services.Err2FARequired) {
		return c.Redirect().To(os.Getenv("FE_URL") + "/login/2fa/" + result.TempToken)
	}
	if err != nil {
		return c.Redirect().To(os.Getenv("FE_URL") + "/login?error=" + url.QueryEscape(err.Error()))
	}
	ctrl.setCookie(c, result.AccessToken, result.RefreshToken)
	if result.User.Role == "admin" {
		return c.Redirect().To(os.Getenv("FE_URL") + "/admin/dashboard")
	}
	return c.Redirect().To(os.Getenv("FE_URL"))
}

func (ctrl *AuthController) DiscordRedirect(c fiber.Ctx) error {
	return c.Redirect().To(ctrl.AuthHTTPService.DiscordRedirectURL())
}

func (ctrl *AuthController) DiscordCallback(c fiber.Ctx) error {
	result, err := ctrl.AuthHTTPService.DiscordCallback(c.Query("code"))
	if errors.Is(err, services.Err2FARequired) {
		return c.Redirect().To(os.Getenv("FE_URL") + "/login/2fa/" + result.TempToken)
	}
	if err != nil {
		return c.Redirect().To(os.Getenv("FE_URL") + "/login?error=" + url.QueryEscape(err.Error()))
	}
	ctrl.setCookie(c, result.AccessToken, result.RefreshToken)
	if result.User.Role == "admin" {
		return c.Redirect().To(os.Getenv("FE_URL") + "/admin/dashboard")
	}
	return c.Redirect().To(os.Getenv("FE_URL"))
}

func (ctrl *AuthController) GoogleOneTap(c fiber.Ctx) error {
	var body struct {
		Credential string `json:"credential"`
	}
	if err := c.Bind().Body(&body); err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: fiber.StatusBadRequest, Message: "Invalid request"}, err)
	}
	result, err := ctrl.AuthHTTPService.GoogleOneTap(body.Credential)
	if errors.Is(err, services.Err2FARequired) {
		return utils.SuccessResponse(c, utils.SuccessResponseFormat{Code: fiber.StatusOK, Message: "2FA diperlukan", Data: fiber.Map{"is_2fa_enabled": true, "temp_token": result.TempToken}})
	}
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: fiber.StatusInternalServerError, Message: "Error handle user google"}, err)
	}
	ctrl.setCookie(c, result.AccessToken, result.RefreshToken)
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{Message: "Success", Data: result.User})
}

func (ctrl *AuthController) FacebookRedirect(c fiber.Ctx) error {
	return c.Redirect().To(ctrl.AuthHTTPService.FacebookRedirectURL())
}

func (ctrl *AuthController) FacebookCallback(c fiber.Ctx) error {
	result, err := ctrl.AuthHTTPService.FacebookCallback(c.Query("code"))
	if errors.Is(err, services.Err2FARequired) {
		return c.Redirect().To(os.Getenv("FE_URL") + "/login/2fa/" + result.TempToken)
	}
	if err != nil {
		return c.Redirect().To(os.Getenv("FE_URL") + "/login?error=" + url.QueryEscape(err.Error()))
	}
	ctrl.setCookie(c, result.AccessToken, result.RefreshToken)
	if result.User.Role == "admin" {
		return c.Redirect().To(os.Getenv("FE_URL") + "/admin/dashboard")
	}
	return c.Redirect().To(os.Getenv("FE_URL"))
}

func (ctrl *AuthController) SteamRedirect(c fiber.Ctx) error {
	return c.Redirect().To(ctrl.AuthHTTPService.SteamRedirectURL())
}

func (ctrl *AuthController) SteamCallback(c fiber.Ctx) error {
	result, err := ctrl.AuthHTTPService.SteamCallback(c.Queries())
	if errors.Is(err, services.Err2FARequired) {
		return c.Redirect().To(os.Getenv("FE_URL") + "/login/2fa/" + result.TempToken)
	}
	if err != nil {
		return c.Redirect().To(os.Getenv("FE_URL") + "/login?error=" + url.QueryEscape(err.Error()))
	}
	ctrl.setCookie(c, result.AccessToken, result.RefreshToken)
	if result.User.Role == "admin" {
		return c.Redirect().To(os.Getenv("FE_URL") + "/admin/dashboard")
	}
	return c.Redirect().To(os.Getenv("FE_URL"))
}

func (ctrl *AuthController) TwitchRedirect(c fiber.Ctx) error {
	return c.Redirect().To(ctrl.AuthHTTPService.TwitchRedirectURL())
}

func (ctrl *AuthController) TwitchCallback(c fiber.Ctx) error {
	result, err := ctrl.AuthHTTPService.TwitchCallback(c.Query("code"))
	if errors.Is(err, services.Err2FARequired) {
		return c.Redirect().To(os.Getenv("FE_URL") + "/login/2fa/" + result.TempToken)
	}
	if err != nil {
		return c.Redirect().To(os.Getenv("FE_URL") + "/login?error=" + url.QueryEscape(err.Error()))
	}
	ctrl.setCookie(c, result.AccessToken, result.RefreshToken)
	if result.User.Role == "admin" {
		return c.Redirect().To(os.Getenv("FE_URL") + "/admin/dashboard")
	}
	return c.Redirect().To(os.Getenv("FE_URL"))
}

func (ctrl *AuthController) Register2FA(c fiber.Ctx) error {
	user := c.Value("user").(jwt.MapClaims)
	secret, qrString, err := ctrl.AuthHTTPService.Register2FA(user["email"].(string))
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: fiber.StatusInternalServerError, Message: "Gagal menyimpan kunci rahasia"}, err)
	}
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{Message: "Success", Data: fiber.Map{"totp_secret": secret, "qr_string": qrString}})
}

func (ctrl *AuthController) Verify2FA(c fiber.Ctx) error {
	user := c.Value("user").(jwt.MapClaims)
	input, err := utils.GetValidatedBody[requests.Verify2FARequest](c)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: fiber.StatusBadRequest, Message: "Gagal parsing request body"}, err)
	}
	if err := ctrl.AuthHTTPService.Verify2FA(user["email"].(string), input.Code, false); err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: fiber.StatusNotFound, Message: "User tidak ditemukan atau token tidak valid"}, err)
	}
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{Message: "Sukses verifikasi 2FA"})
}

func (ctrl *AuthController) Disable2FA(c fiber.Ctx) error {
	user := c.Value("user").(jwt.MapClaims)
	if err := ctrl.AuthHTTPService.Disable2FA(user["id"].(string)); err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: fiber.StatusInternalServerError, Message: "Gagal menghapus 2FA"}, err)
	}
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{Message: "Sukses menghapus 2FA"})
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
	if err := ctrl.AuthHTTPService.SendOTP(input.TargetType, target, input.Subject, userName); err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: fiber.StatusInternalServerError, Message: "Gagal mengirim OTP"}, err)
	}
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{Message: "Success"})
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
	if err := ctrl.AuthHTTPService.VerifyOTP(target, input.Subject, input.Otp); err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: fiber.StatusNotFound, Message: "Kode OTP tidak valid"}, err)
	}
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{Message: "Success"})
}

func (ctrl *AuthController) setCookie(c fiber.Ctx, accessToken string, refreshToken string) {
	accessCookie := fiber.Cookie{Name: "access_token_" + os.Getenv("APP_ENV"), Value: accessToken, Expires: time.Now().Add(time.Hour), HTTPOnly: true, Domain: os.Getenv("FE_DOMAIN"), SameSite: fiber.CookieSameSiteNoneMode, Secure: true}
	refreshCookie := fiber.Cookie{Name: "refresh_token_" + os.Getenv("APP_ENV"), Value: refreshToken, Expires: time.Now().Add(time.Hour * 24), HTTPOnly: true, Domain: os.Getenv("FE_DOMAIN"), SameSite: fiber.CookieSameSiteNoneMode, Secure: true}
	c.Cookie(&accessCookie)
	c.Cookie(&refreshCookie)
}

func (ctrl *AuthController) clearCookie(c fiber.Ctx) {
	c.Cookie(&fiber.Cookie{Name: "access_token_" + os.Getenv("APP_ENV"), Value: "", Expires: time.Now().Add(-time.Hour), HTTPOnly: true, Domain: os.Getenv("FE_DOMAIN"), SameSite: fiber.CookieSameSiteNoneMode, Secure: true})
	c.Cookie(&fiber.Cookie{Name: "refresh_token_" + os.Getenv("APP_ENV"), Value: "", Expires: time.Now().Add(-time.Hour), HTTPOnly: true, Domain: os.Getenv("FE_DOMAIN"), SameSite: fiber.CookieSameSiteNoneMode, Secure: true})
}
