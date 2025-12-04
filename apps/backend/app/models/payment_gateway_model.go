package models

import (
	"time"

	"gorm.io/gorm"
)

type PaymentGateway struct {
	ID             string    `gorm:"primaryKey;size:7" json:"id"`
	Order          int       `gorm:"not null;index" json:"order"`
	Name           string    `gorm:"not null;size:100" json:"name"`
	Metadata       *string   `gorm:"type:text" json:"metadata"`
	Description    *string   `gorm:"type:text" json:"description"`
	IsActive       bool      `gorm:"default:true;not null;index" json:"is_active"`
	CreatedAt      time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt      time.Time `gorm:"not null" json:"updated_at"`
	PaymentMethods []PaymentMethod
}

func (t *PaymentGateway) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == "" {
		t.ID = GenerateID(7)
	}

	return
}
