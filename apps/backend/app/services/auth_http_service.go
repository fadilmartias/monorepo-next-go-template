package services

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/fadilmartias/dilz_code/apps/backend/app/models"
	"github.com/fadilmartias/dilz_code/apps/backend/app/responses"
	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"
	"github.com/fadilmartias/dilz_code/apps/backend/config"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
)

type AuthHTTPService struct {
	Auth         *AuthService
	UserService  *UserService
	EmailService *EmailService
	DB           *gorm.DB
	Redis        *redis.Client
}

type OAuthResult struct {
	User         *models.User
	AccessToken  string
	RefreshToken string
	TempToken    string
	RedirectURL  string
	Is2FA        bool
}

func NewAuthHTTPService(auth *AuthService, userService *UserService, emailService *EmailService, db *gorm.DB, redisClient *redis.Client) *AuthHTTPService {
	return &AuthHTTPService{
		Auth:         auth,
		UserService:  userService,
		EmailService: emailService,
		DB:           db,
		Redis:        redisClient,
	}
}

func (s *AuthHTTPService) Me(userClaims jwt.MapClaims) (*responses.UserResponse, error) {
	userID, _ := userClaims["id"].(string)
	if userID == "" {
		return nil, errors.New("invalid token")
	}
	return s.UserService.FindByID(userID)
}

func (s *AuthHTTPService) Logout(userClaims jwt.MapClaims) error {
	email, _ := userClaims["email"].(string)
	if email == "" {
		return errors.New("invalid token")
	}

	var userDB models.User
	if err := s.DB.Where("email = ?", email).First(&userDB).Error; err != nil {
		return err
	}

	userDB.RefreshToken = nil
	return s.DB.Save(&userDB).Error
}

func (s *AuthHTTPService) ForgotPassword(email string) (string, error) {
	var user models.User
	if err := s.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return "", err
	}

	var passwordResetToken models.PasswordResetToken
	if err := s.DB.Where("email = ?", email).Where("expired_at > ?", time.Now()).First(&passwordResetToken).Error; err == nil {
		return "", nil
	}

	token := utils.GenerateRandomToken(32)
	passwordResetToken.Email = user.Email
	passwordResetToken.Token = utils.HashToken(token)
	passwordResetToken.ExpiredAt = time.Now().Add(time.Minute * 5)
	if err := s.DB.Save(&passwordResetToken).Error; err != nil {
		return "", err
	}

	if err := s.EmailService.SendEmailResetPassword(user.Email, token, user.Name); err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthHTTPService) SendEmailVerification(userClaims jwt.MapClaims) (string, error) {
	userID, _ := userClaims["id"].(string)
	if userID == "" {
		return "", errors.New("invalid token")
	}

	var user models.User
	if err := s.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return "", err
	}
	if user.EmailVerifiedAt != nil {
		return "", nil
	}

	key := fmt.Sprintf("email_verification_token:%s", user.Email)
	if _, err := s.Redis.Get(context.Background(), key).Result(); err == nil {
		return "", errors.New("email verification already sent")
	}

	token, err := utils.GenerateToken(map[string]any{"email": user.Email}, time.Minute*5)
	if err != nil {
		return "", err
	}

	if err := s.Redis.Set(context.Background(), key, token, time.Minute*5).Err(); err != nil {
		return "", err
	}

	if err := s.EmailService.SendEmailVerification(user.Email, token, user.Name); err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthHTTPService) VerifyEmail(token string) (*models.User, error) {
	if token == "" {
		return nil, errors.New("token invalid")
	}

	claims, err := utils.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	var user models.User
	if err := s.DB.Where("email = ?", claims["email"]).First(&user).Error; err != nil {
		return nil, err
	}
	if user.EmailVerifiedAt != nil {
		return &user, nil
	}

	now := time.Now()
	user.EmailVerifiedAt = &now
	if err := s.DB.Save(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *AuthHTTPService) RefreshAccessToken(refreshToken string) (string, string, *models.User, error) {
	claims, err := utils.ValidateToken(refreshToken)
	if err != nil {
		return "", "", nil, err
	}

	var user models.User
	if err := s.DB.Where("email = ?", claims["email"]).First(&user).Error; err != nil {
		return "", "", nil, err
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
		return "", "", nil, err
	}

	newRefreshToken, err := utils.GenerateToken(map[string]any{
		"id":             user.ID,
		"name":           user.Name,
		"tenant_id":      user.TenantID,
		"phone":          user.Phone,
		"email":          user.Email,
		"role":           user.Role,
		"is_2fa_enabled": user.TOTPSecret != nil,
	}, time.Hour*24)
	if err != nil {
		return "", "", nil, err
	}

	return accessToken, newRefreshToken, &user, nil
}

func (s *AuthHTTPService) GoogleRedirectURL() string {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	redirectURI := os.Getenv("APP_URL") + "/v1/auth/google/callback"
	return fmt.Sprintf(
		"https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=email profile",
		clientID, redirectURI,
	)
}

func (s *AuthHTTPService) GoogleCallback(code string) (*OAuthResult, error) {
	if code == "" {
		return nil, errors.New("code not found")
	}

	googleConfig := config.LoadGoogleConfig()
	redirectURI := os.Getenv("APP_URL") + "/v1/auth/google/callback"

	resp, err := utils.Http().WithFormData(map[string]string{
		"code":          code,
		"client_id":     googleConfig.ClientID,
		"client_secret": googleConfig.ClientSecret,
		"redirect_uri":  redirectURI,
		"grant_type":    "authorization_code",
	}).Post("https://oauth2.googleapis.com/token")
	if err != nil {
		return nil, err
	}

	accessTokenGoogle := gjson.GetBytes(resp.Body(), "access_token").String()
	userResp, err := utils.Http().WithAuthToken(accessTokenGoogle).Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}

	id := gjson.GetBytes(userResp.Body(), "id").String()
	email := gjson.GetBytes(userResp.Body(), "email").String()
	name := gjson.GetBytes(userResp.Body(), "name").String()
	user, accessToken, refreshToken, tempToken, _, err := s.Auth.HandleOAuth(context.Background(), "google", id, email, name)
	if err != nil {
		return &OAuthResult{TempToken: tempToken}, err
	}
	return &OAuthResult{User: user, AccessToken: accessToken, RefreshToken: refreshToken, TempToken: tempToken}, nil
}

func (s *AuthHTTPService) DiscordRedirectURL() string {
	discordConfig := config.LoadDiscordConfig()
	redirectURI := os.Getenv("APP_URL") + "/v1/auth/discord/callback"
	return fmt.Sprintf(
		"https://discord.com/api/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=identify+email",
		discordConfig.ClientID, redirectURI,
	)
}

func (s *AuthHTTPService) DiscordCallback(code string) (*OAuthResult, error) {
	if code == "" {
		return nil, errors.New("missing code")
	}
	discordConfig := config.LoadDiscordConfig()
	appConfig := config.LoadAppConfig()
	redirectURI := appConfig.BaseURL + "/v1/auth/discord/callback"

	resp, err := utils.Http().WithFormData(map[string]string{
		"code":          code,
		"client_id":     discordConfig.ClientID,
		"client_secret": discordConfig.ClientSecret,
		"redirect_uri":  redirectURI,
		"grant_type":    "authorization_code",
	}).Post("https://discord.com/api/oauth2/token")
	if err != nil {
		return nil, err
	}

	accessTokenDiscord := gjson.GetBytes(resp.Body(), "access_token").String()
	userResp, err := utils.Http().WithAuthToken(accessTokenDiscord).Get("https://discord.com/api/users/@me")
	if err != nil {
		return nil, err
	}

	discordID := gjson.GetBytes(userResp.Body(), "id").String()
	username := gjson.GetBytes(userResp.Body(), "username").String()
	email := gjson.GetBytes(userResp.Body(), "email").String()
	user, accessToken, refreshToken, tempToken, _, err := s.Auth.HandleOAuth(context.Background(), "discord", discordID, email, username)
	if err != nil {
		return &OAuthResult{TempToken: tempToken}, err
	}
	return &OAuthResult{User: user, AccessToken: accessToken, RefreshToken: refreshToken, TempToken: tempToken}, nil
}

func (s *AuthHTTPService) GoogleOneTap(credential string) (*OAuthResult, error) {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	resp, err := utils.Http().WithQuery(map[string]string{"id_token": credential}).Get("https://oauth2.googleapis.com/tokeninfo")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, errors.New("invalid id token")
	}

	email := gjson.GetBytes(resp.Body(), "email").String()
	name := gjson.GetBytes(resp.Body(), "name").String()
	aud := gjson.GetBytes(resp.Body(), "aud").String()
	id := gjson.GetBytes(resp.Body(), "id").String()
	if aud != clientID {
		return nil, errors.New("error credential")
	}
	user, accessToken, refreshToken, tempToken, _, err := s.Auth.HandleOAuth(context.Background(), "google", id, email, name)
	if err != nil {
		return &OAuthResult{TempToken: tempToken}, err
	}
	return &OAuthResult{User: user, AccessToken: accessToken, RefreshToken: refreshToken, TempToken: tempToken}, nil
}

func (s *AuthHTTPService) FacebookRedirectURL() string {
	facebookConfig := config.LoadFacebookConfig()
	appConfig := config.LoadAppConfig()
	redirectURI := appConfig.BaseURL + "/v1/auth/facebook/callback"
	return fmt.Sprintf(
		"https://www.facebook.com/v23.0/dialog/oauth?client_id=%s&redirect_uri=%s&scope=email,public_profile&response_type=code&state=%s",
		facebookConfig.ClientID, url.QueryEscape(redirectURI), "randomstate",
	)
}

func (s *AuthHTTPService) FacebookCallback(code string) (*OAuthResult, error) {
	if code == "" {
		return nil, errors.New("missing code")
	}
	facebookConfig := config.LoadFacebookConfig()
	appConfig := config.LoadAppConfig()
	redirectURI := appConfig.BaseURL + "/v1/auth/facebook/callback"
	tokenURL := fmt.Sprintf(
		"https://graph.facebook.com/v23.0/oauth/access_token?client_id=%s&redirect_uri=%s&client_secret=%s&code=%s",
		facebookConfig.ClientID, url.QueryEscape(redirectURI), facebookConfig.ClientSecret, code,
	)
	resp, err := utils.Http().Get(tokenURL)
	if err != nil {
		return nil, err
	}
	accessToken := gjson.GetBytes(resp.Body(), "access_token").String()
	userInfoURL := fmt.Sprintf("https://graph.facebook.com/me?fields=id,name,email&access_token=%s", accessToken)
	userResp, err := utils.Http().Get(userInfoURL)
	if err != nil {
		return nil, err
	}
	id := gjson.GetBytes(userResp.Body(), "id").String()
	name := gjson.GetBytes(userResp.Body(), "name").String()
	email := gjson.GetBytes(userResp.Body(), "email").String()
	user, accessTokenOut, refreshToken, tempToken, _, err := s.Auth.HandleOAuth(context.Background(), "facebook", id, email, name)
	if err != nil {
		return &OAuthResult{TempToken: tempToken}, err
	}
	return &OAuthResult{User: user, AccessToken: accessTokenOut, RefreshToken: refreshToken, TempToken: tempToken}, nil
}

func (s *AuthHTTPService) SteamRedirectURL() string {
	appConfig := config.LoadAppConfig()
	redirectURI := appConfig.BaseURL + "/v1/auth/steam/callback"
	return fmt.Sprintf(
		"https://steamcommunity.com/openid/login?openid.return_to=%s&openid.realm=%s",
		url.QueryEscape(redirectURI), url.QueryEscape(appConfig.BaseURL),
	)
}

func (s *AuthHTTPService) SteamCallback(query map[string]string) (*OAuthResult, error) {
	form := url.Values{}
	for k, v := range query {
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

	resp, err := utils.Http().WithFormData(formData).Post("https://steamcommunity.com/openid/login")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 || !strings.Contains(string(resp.Body()), "is_valid:true") {
		return nil, errors.New("invalid steam login")
	}

	claimedID := query["openid.claimed_id"]
	if claimedID == "" {
		return nil, errors.New("missing claimed_id")
	}
	steamID := claimedID[strings.LastIndex(claimedID, "/")+1:]
	name := "SteamUser_" + steamID
	if apiKey := os.Getenv("STEAM_API_KEY"); apiKey != "" {
		profileURL := fmt.Sprintf("https://api.steampowered.com/ISteamUser/GetPlayerSummaries/v2/?key=%s&steamids=%s", apiKey, steamID)
		profileResp, err := utils.Http().Get(profileURL)
		if err == nil && profileResp.StatusCode() == 200 {
			name = gjson.GetBytes(profileResp.Body(), "response.players.0.personaname").String()
		}
	}
	user, accessToken, refreshToken, tempToken, _, err := s.Auth.HandleOAuth(context.Background(), "steam", steamID, "", name)
	if err != nil {
		return &OAuthResult{TempToken: tempToken}, err
	}
	return &OAuthResult{User: user, AccessToken: accessToken, RefreshToken: refreshToken, TempToken: tempToken}, nil
}

func (s *AuthHTTPService) TwitchRedirectURL() string {
	twitchConfig := config.LoadTwitchConfig()
	appConfig := config.LoadAppConfig()
	redirectURI := appConfig.BaseURL + "/v1/auth/twitch/callback"
	return fmt.Sprintf(
		"https://id.twitch.tv/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=user:read:email",
		twitchConfig.ClientID, url.QueryEscape(redirectURI),
	)
}

func (s *AuthHTTPService) TwitchCallback(code string) (*OAuthResult, error) {
	if code == "" {
		return nil, errors.New("missing code")
	}
	twitchConfig := config.LoadTwitchConfig()
	appConfig := config.LoadAppConfig()
	redirectURI := appConfig.BaseURL + "/v1/auth/twitch/callback"

	tokenResp, err := utils.Http().WithHeader("Content-Type", "application/x-www-form-urlencoded").WithFormData(map[string]string{
		"client_id":     twitchConfig.ClientID,
		"client_secret": twitchConfig.ClientSecret,
		"code":          code,
		"grant_type":    "authorization_code",
		"redirect_uri":  redirectURI,
	}).Post("https://id.twitch.tv/oauth2/token")
	if err != nil {
		return nil, err
	}
	accessToken := gjson.GetBytes(tokenResp.Body(), "access_token").String()
	userResp, err := utils.Http().WithHeaders(map[string]string{
		"Authorization": "Bearer " + accessToken,
		"Client-Id":     twitchConfig.ClientID,
	}).Get("https://api.twitch.tv/helix/users")
	if err != nil {
		return nil, err
	}

	id := gjson.GetBytes(userResp.Body(), "data.0.id").String()
	email := gjson.GetBytes(userResp.Body(), "data.0.email").String()
	name := gjson.GetBytes(userResp.Body(), "data.0.display_name").String()
	user, accessTokenOut, refreshToken, tempToken, _, err := s.Auth.HandleOAuth(context.Background(), "twitch", id, email, name)
	if err != nil {
		return &OAuthResult{TempToken: tempToken}, err
	}
	return &OAuthResult{User: user, AccessToken: accessTokenOut, RefreshToken: refreshToken, TempToken: tempToken}, nil
}

func (s *AuthHTTPService) Register2FA(email string) (string, string, error) {
	return s.Auth.Register2FA(email)
}

func (s *AuthHTTPService) Verify2FA(email string, code string, isLogin bool) error {
	return s.Auth.VerifyUserTOTPSecret(email, code, isLogin)
}

func (s *AuthHTTPService) Disable2FA(id string) error {
	return s.Auth.Disable2FA(id)
}

func (s *AuthHTTPService) SendOTP(targetType string, target string, subject string, name string) error {
	return s.Auth.SendOTP(context.Background(), targetType, target, subject, name)
}

func (s *AuthHTTPService) VerifyOTP(target string, subject string, otp string) error {
	return s.Auth.VerifyOTP(context.Background(), target, subject, otp)
}

func (s *AuthHTTPService) ForgotPasswordForUser(email string) (string, error) {
	var user models.User
	if err := s.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return "", err
	}

	var passwordResetToken models.PasswordResetToken
	if err := s.DB.Where("email = ?", email).Where("expired_at > ?", time.Now()).First(&passwordResetToken).Error; err == nil {
		return "", nil
	}

	token := utils.GenerateRandomToken(32)
	passwordResetToken.Email = user.Email
	passwordResetToken.Token = utils.HashToken(token)
	passwordResetToken.ExpiredAt = time.Now().Add(time.Minute * 5)
	if err := s.DB.Save(&passwordResetToken).Error; err != nil {
		return "", err
	}
	if err := s.EmailService.SendEmailResetPassword(user.Email, token, user.Name); err != nil {
		return "", err
	}
	return token, nil
}

func (s *AuthHTTPService) ResetPassword(token, password string) error {
	tokenHash := utils.HashToken(token)
	var passwordResetToken models.PasswordResetToken
	if err := s.DB.Where("token = ?", tokenHash).First(&passwordResetToken).Error; err != nil {
		return err
	}
	if passwordResetToken.UsedAt != nil {
		return errors.New("password reset token already used")
	}
	var user models.User
	if err := s.DB.Where("email = ?", passwordResetToken.Email).First(&user).Error; err != nil {
		return err
	}
	if err := user.HashPassword(password); err != nil {
		return err
	}
	if err := s.DB.Save(&user).Error; err != nil {
		return err
	}
	now := time.Now()
	passwordResetToken.UsedAt = &now
	return s.DB.Save(&passwordResetToken).Error
}

func (s *AuthHTTPService) SendEmailVerificationToken(userClaims jwt.MapClaims) (string, error) {
	return s.SendEmailVerification(userClaims)
}

func (s *AuthHTTPService) RefreshClaimsToResponse(user *models.User) responses.UserResponse {
	return responses.UserResponse{
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
}
