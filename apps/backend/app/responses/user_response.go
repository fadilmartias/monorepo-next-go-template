package responses

import (
	"time"

	"github.com/fadilmartias/dilz_code/apps/backend/app/models"
)

type UserResponse struct {
	ID              *string          `json:"id,omitempty"`
	TenantID        *string          `json:"tenant_id,omitempty"`
	Name            *string          `json:"name,omitempty"`
	Email           *string          `json:"email,omitempty"`
	Phone           *string          `json:"phone,omitempty"`
	Role            *models.UserRole `json:"role,omitempty"`
	EmailVerifiedAt *time.Time       `json:"email_verified_at,omitempty"`
	CreatedAt       *time.Time       `json:"created_at,omitempty"`
	UpdatedAt       *time.Time       `json:"updated_at,omitempty"`
	DeletedAt       *time.Time       `json:"deleted_at,omitempty"`
	Is2FAEnabled    bool             `json:"is_2fa_enabled,omitempty"`
	Level           int              `json:"level"`
	Exp             int              `json:"exp"`
	TotalSpent      float64          `json:"total_spent"`
	TotalExp        int              `json:"total_exp"`
	ExpForLevel     int              `json:"exp_for_level"`
	Point           float64          `json:"point,omitempty"`
	ReferralCode    string           `json:"referral_code,omitempty"`
}
