package models

type UserPasskey struct {
	BaseModelWithoutUpdatedAt
	UserID       string `gorm:"type:char(36);index" json:"user_id"`
	CredentialID string `gorm:"not null;index" json:"credential_id"`
	PublicKey    string `gorm:"not null;type:text" json:"public_key"`
	SignCount    int    `gorm:"not null;default:0" json:"sign_count"`
	IsActive     bool   `gorm:"not null;default:true" json:"is_active"`
}
