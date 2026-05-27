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

// VerifyRecaptcha memvalidasi token Google reCAPTCHA
// @Summary Validasi reCAPTCHA
// @Description Memvalidasi token reCAPTCHA yang dikirim dari sisi klien sebelum memproses form sensitif.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param payload body ValidateRecaptchaInput true "Data token dan action reCAPTCHA"
// @Success 200 {object} utils.OrderedSuccessResponse "Validasi recaptcha berhasil"
// @Failure 400 {object} utils.OrderedErrorResponse "Token dan action harus diisi"
// @Failure 401 {object} utils.OrderedErrorResponse "Validasi recaptcha gagal"
// @Router /v1/auth/verify-recaptcha [post]
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

// Me mengambil data profil user yang sedang login
// @Summary Dapatkan Profil User (Me)
// @Description Mengambil data profil pengguna berdasarkan token otentikasi yang aktif (via Cookie/Bearer).
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.OrderedSuccessResponse "User retrieved successfully"
// @Failure 401 {object} utils.OrderedErrorResponse "Invalid token"
// @Failure 404 {object} utils.OrderedErrorResponse "User tidak ditemukan"
// @Router /v1/auth/me [get]
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

// Logout mengeluarkan pengguna dan menghapus cookie
// @Summary Logout User
// @Description Mengakhiri sesi pengguna dengan menghapus token otentikasi dari cookie.
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.OrderedSuccessResponse "Logged out successfully"
// @Failure 401 {object} utils.OrderedErrorResponse "Invalid token"
// @Failure 500 {object} utils.OrderedErrorResponse "Failed to logout user"
// @Router /v1/auth/logout [post]
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

// Register mendaftarkan akun baru
// @Summary Registrasi Akun Baru
// @Description Mendaftarkan pengguna baru ke dalam sistem.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param payload body requests.RegisterInput true "Data pendaftaran pengguna"
// @Success 201 {object} utils.OrderedSuccessResponse "User created successfully"
// @Failure 400 {object} utils.OrderedErrorResponse "Bad Request (Validasi gagal)"
// @Failure 500 {object} utils.OrderedErrorResponse "Internal Server Error"
// @Router /v1/auth/register [post]
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

// Login memproses masuk pengguna
// @Summary Login User
// @Description Memproses kredensial (email/username) dan password. Akan mengatur cookie `access_token` dan `refresh_token` jika berhasil. Mengembalikan flag khusus jika akun mengaktifkan 2FA.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param payload body requests.LoginInput true "Kredensial login"
// @Success 200 {object} utils.OrderedSuccessResponse "Login berhasil atau membutuhkan 2FA"
// @Failure 400 {object} utils.OrderedErrorResponse "Validasi request gagal"
// @Failure 401 {object} utils.OrderedErrorResponse "Email atau password salah"
// @Router /v1/auth/login [post]
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

// Login2FA memproses login menggunakan kode 2FA
// @Summary Login via 2FA
// @Description Menyelesaikan proses login dengan menggunakan `temp_token` dari tahap login pertama dan kode otentikator.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param payload body requests.Login2FARequest true "Temp token dan kode 2FA"
// @Success 200 {object} utils.OrderedSuccessResponse "Login berhasil"
// @Failure 400 {object} utils.OrderedErrorResponse "Validasi request gagal"
// @Failure 401 {object} utils.OrderedErrorResponse "Kode 2FA salah atau token kedaluwarsa"
// @Router /v1/auth/login-2fa [post]
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

// ForgotPassword memicu pengiriman email reset password
// @Summary Lupa Password
// @Description Meminta link token reset password dikirimkan ke email yang terdaftar.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param payload body requests.ForgotPasswordInput true "Email pengguna"
// @Success 200 {object} utils.OrderedSuccessResponse "Token reset password berhasil dikirim"
// @Failure 400 {object} utils.OrderedErrorResponse "Validasi request gagal"
// @Failure 500 {object} utils.OrderedErrorResponse "Gagal memproses forgot password"
// @Router /v1/auth/forgot-password [post]
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

// ResetPassword mengubah password menggunakan token
// @Summary Reset Password
// @Description Mengatur ulang kata sandi dengan memberikan token yang valid (didapat dari email) beserta kata sandi baru.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param payload body requests.ResetPasswordInput true "Token dan password baru"
// @Success 200 {object} utils.OrderedSuccessResponse "Password reset berhasil"
// @Failure 400 {object} utils.OrderedErrorResponse "Validasi request gagal"
// @Failure 500 {object} utils.OrderedErrorResponse "Gagal mereset password atau token tidak valid"
// @Router /v1/auth/reset-password [post]
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

// SendEmailVerification mengirim ulang email verifikasi
// @Summary Kirim Ulang Email Verifikasi
// @Description Mengirimkan kembali tautan verifikasi alamat email ke pengguna yang sedang login.
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.OrderedSuccessResponse "Email verifikasi berhasil dikirim"
// @Failure 401 {object} utils.OrderedErrorResponse "Invalid token"
// @Failure 500 {object} utils.OrderedErrorResponse "Gagal mengirim email verifikasi"
// @Router /v1/auth/send-email-verification [post]
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

// VerifyEmail memverifikasi alamat email menggunakan token
// @Summary Verifikasi Email
// @Description Memverifikasi bahwa email milik pengguna menggunakan token rahasia yang dikirimkan via email.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param payload body requests.VerifyEmailInput true "Token verifikasi email"
// @Success 200 {object} utils.OrderedSuccessResponse "Email verified successfully"
// @Failure 400 {object} utils.OrderedErrorResponse "Validasi request gagal"
// @Failure 401 {object} utils.OrderedErrorResponse "Token tidak valid"
// @Router /v1/auth/verify-email [post]
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

// RefreshAccessToken memperbarui access token via refresh token
// @Summary Refresh Token
// @Description Menghasilkan `access_token` yang baru menggunakan `refresh_token` yang valid dari Cookie.
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} utils.OrderedSuccessResponse "Token refreshed successfully"
// @Failure 400 {object} utils.OrderedErrorResponse "Token tidak ditemukan atau tidak valid"
// @Router /v1/auth/refresh-access-token [post]
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

// GoogleRedirect me-redirect user ke halaman login Google
// @Summary Login Google (Redirect)
// @Description Mengarahkan (redirect) pengguna ke halaman otorisasi OAuth Google.
// @Tags OAuth
// @Router /v1/auth/google/redirect [get]
func (ctrl *AuthController) GoogleRedirect(c fiber.Ctx) error {
	return c.Redirect().To(ctrl.AuthHTTPService.GoogleRedirectURL())
}

// GoogleCallback menangani callback dari Google OAuth
// @Summary Callback Google OAuth
// @Description Menangani respons balikan (code) dari Google setelah login sukses.
// @Tags OAuth
// @Param code query string true "Kode otorisasi dari Google"
// @Router /v1/auth/google/callback [get]
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

// DiscordRedirect me-redirect user ke halaman login Discord
// @Summary Login Discord (Redirect)
// @Description Mengarahkan pengguna ke halaman otorisasi OAuth Discord.
// @Tags OAuth
// @Router /v1/auth/discord/redirect [get]
func (ctrl *AuthController) DiscordRedirect(c fiber.Ctx) error {
	return c.Redirect().To(ctrl.AuthHTTPService.DiscordRedirectURL())
}

// DiscordCallback menangani callback dari Discord OAuth
// @Summary Callback Discord OAuth
// @Description Menangani respons balikan dari Discord.
// @Tags OAuth
// @Param code query string true "Kode otorisasi dari Discord"
// @Router /v1/auth/discord/callback [get]
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

// GoogleOneTap memproses login menggunakan Google One Tap
// @Summary Google One Tap Login
// @Description Memvalidasi kredensial Google One Tap JWT yang dikirim dari klien.
// @Tags OAuth
// @Accept json
// @Produce json
// @Param payload body map[string]string true "Format JSON dengan property 'credential'"
// @Success 200 {object} utils.OrderedSuccessResponse "Berhasil login atau membutuhkan 2FA"
// @Failure 400 {object} utils.OrderedErrorResponse "Invalid request"
// @Failure 500 {object} utils.OrderedErrorResponse "Error handle user google"
// @Router /v1/auth/google/one-tap [post]
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

// FacebookRedirect me-redirect user ke halaman login Facebook
// @Summary Login Facebook (Redirect)
// @Description Mengarahkan pengguna ke halaman otorisasi OAuth Facebook.
// @Tags OAuth
// @Router /v1/auth/facebook/redirect [get]
func (ctrl *AuthController) FacebookRedirect(c fiber.Ctx) error {
	return c.Redirect().To(ctrl.AuthHTTPService.FacebookRedirectURL())
}

// FacebookCallback menangani callback dari Facebook OAuth
// @Summary Callback Facebook OAuth
// @Description Menangani respons balikan dari Facebook.
// @Tags OAuth
// @Param code query string true "Kode otorisasi dari Facebook"
// @Router /v1/auth/facebook/callback [get]
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

// SteamRedirect me-redirect user ke halaman login Steam
// @Summary Login Steam (Redirect)
// @Description Mengarahkan pengguna ke halaman OpenID Steam.
// @Tags OAuth
// @Router /v1/auth/steam/redirect [get]
func (ctrl *AuthController) SteamRedirect(c fiber.Ctx) error {
	return c.Redirect().To(ctrl.AuthHTTPService.SteamRedirectURL())
}

// SteamCallback menangani callback dari Steam OpenID
// @Summary Callback Steam OpenID
// @Description Menangani respons balikan (query params) dari Steam.
// @Tags OAuth
// @Router /v1/auth/steam/callback [get]
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

// TwitchRedirect me-redirect user ke halaman login Twitch
// @Summary Login Twitch (Redirect)
// @Description Mengarahkan pengguna ke halaman otorisasi OAuth Twitch.
// @Tags OAuth
// @Router /v1/auth/twitch/redirect [get]
func (ctrl *AuthController) TwitchRedirect(c fiber.Ctx) error {
	return c.Redirect().To(ctrl.AuthHTTPService.TwitchRedirectURL())
}

// TwitchCallback menangani callback dari Twitch OAuth
// @Summary Callback Twitch OAuth
// @Description Menangani respons balikan dari Twitch.
// @Tags OAuth
// @Param code query string true "Kode otorisasi dari Twitch"
// @Router /v1/auth/twitch/callback [get]
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

// Register2FA menginisialisasi pendaftaran 2FA untuk user
// @Summary Setup 2FA (TOTP)
// @Description Meng-generate secret TOTP dan URL QR Code untuk dipindai oleh Google Authenticator / Authy.
// @Tags Two-Factor Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.OrderedSuccessResponse "Berhasil generate TOTP secret"
// @Failure 500 {object} utils.OrderedErrorResponse "Gagal menyimpan kunci rahasia"
// @Router /v1/auth/register-2fa [post]
func (ctrl *AuthController) Register2FA(c fiber.Ctx) error {
	user := c.Value("user").(jwt.MapClaims)
	secret, qrString, err := ctrl.AuthHTTPService.Register2FA(user["email"].(string))
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: fiber.StatusInternalServerError, Message: "Gagal menyimpan kunci rahasia"}, err)
	}
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{Message: "Success", Data: fiber.Map{"totp_secret": secret, "qr_string": qrString}})
}

// Verify2FA memverifikasi setup 2FA untuk pertama kali
// @Summary Verifikasi Aktivasi 2FA
// @Description Mengaktifkan 2FA di akun dengan memvalidasi kode PIN pertama yang dimasukkan pengguna.
// @Tags Two-Factor Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body requests.Verify2FARequest true "Kode OTP dari Authenticator"
// @Success 200 {object} utils.OrderedSuccessResponse "Sukses verifikasi 2FA"
// @Failure 400 {object} utils.OrderedErrorResponse "Gagal parsing request body"
// @Failure 404 {object} utils.OrderedErrorResponse "User tidak ditemukan atau token tidak valid"
// @Router /v1/auth/verify-2fa [post]
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

// Disable2FA menonaktifkan fitur 2FA pada akun user
// @Summary Nonaktifkan 2FA
// @Description Mematikan fitur Two-Factor Authentication untuk akun pengguna yang sedang login.
// @Tags Two-Factor Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.OrderedSuccessResponse "Sukses menghapus 2FA"
// @Failure 500 {object} utils.OrderedErrorResponse "Gagal menghapus 2FA"
// @Router /v1/auth/disable-2fa [post]
func (ctrl *AuthController) Disable2FA(c fiber.Ctx) error {
	user := c.Value("user").(jwt.MapClaims)
	if err := ctrl.AuthHTTPService.Disable2FA(user["id"].(string)); err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: fiber.StatusInternalServerError, Message: "Gagal menghapus 2FA"}, err)
	}
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{Message: "Sukses menghapus 2FA"})
}

// SendOTP memicu pengiriman OTP via Email/WA
// @Summary Kirim Kode OTP
// @Description Mengirimkan kode OTP ke email atau WhatsApp. Bisa digunakan oleh pengguna login maupun tamu (guest).
// @Tags OTP
// @Accept json
// @Produce json
// @Param payload body requests.SendOTPRequest true "Data target tujuan pengiriman OTP"
// @Success 200 {object} utils.OrderedSuccessResponse "Sukses mengirim OTP"
// @Failure 400 {object} utils.OrderedErrorResponse "Validasi request gagal"
// @Failure 500 {object} utils.OrderedErrorResponse "Gagal mengirim OTP"
// @Router /v1/auth/send-otp [post]
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

// VerifyOTP memvalidasi kode OTP yang dikirim user
// @Summary Verifikasi Kode OTP
// @Description Memverifikasi input kode OTP yang diterima pengguna.
// @Tags OTP
// @Accept json
// @Produce json
// @Param payload body requests.VerifyOTPRequest true "Data target dan kode OTP"
// @Success 200 {object} utils.OrderedSuccessResponse "Sukses memverifikasi OTP"
// @Failure 400 {object} utils.OrderedErrorResponse "Validasi request gagal"
// @Failure 404 {object} utils.OrderedErrorResponse "Kode OTP tidak valid"
// @Router /v1/auth/verify-otp [post]
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
