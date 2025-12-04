package models

import (
	"time"

	"gorm.io/gorm"
)

type UserPasskey struct {
	ID           string    `gorm:"primarykey;size:7" json:"id"`
	UserID       string    `gorm:"size:7;index;not null" json:"user_id"`
	CredentialID string    `gorm:"not null;index" json:"credential_id"`
	PublicKey    string    `gorm:"size:7;index;not null" json:"public_key"`
	SignCount    int       `gorm:"not null;default:0" json:"sign_count"`
	IsActive     bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt    time.Time `gorm:"not null" json:"created_at"`
}

func (t *UserPasskey) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == "" {
		t.ID = GenerateID(7)
	}

	return
}
