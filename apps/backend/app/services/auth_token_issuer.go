package services

import (
	"time"

	"github.com/fadilmartias/dilz_code/apps/backend/app/models"
	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"
)

type AuthTokenIssuer struct{}

func NewAuthTokenIssuer() *AuthTokenIssuer {
	return &AuthTokenIssuer{}
}

func (i *AuthTokenIssuer) IssueTempToken(email string) (string, error) {
	return utils.GenerateToken(map[string]any{
		"email": email,
	}, time.Minute*5)
}

func (i *AuthTokenIssuer) Claims(user *models.User, is2FAEnabled bool) map[string]any {
	return map[string]any{
		"id":             user.ID,
		"name":           user.Name,
		"tenant_id":      user.TenantID,
		"phone":          user.Phone,
		"email":          user.Email,
		"role":           user.Role,
		"is_2fa_enabled": is2FAEnabled,
	}
}

func (i *AuthTokenIssuer) IssueTokenPair(user *models.User, is2FAEnabled bool) (string, string, error) {
	claims := i.Claims(user, is2FAEnabled)
	accessToken, err := utils.GenerateToken(claims, time.Hour*1)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := utils.GenerateToken(claims, time.Hour*24)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
