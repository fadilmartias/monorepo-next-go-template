package models

import (
	"time"

	"gorm.io/gorm"
)

type Setting struct {
	ID          string    `gorm:"primarykey;size:7" json:"id"`
	Key         string    `gorm:"not null;uniqueIndex;size:255" json:"key"`
	Value       string    `gorm:"not null;type:text" json:"value"`
	Type        string    `gorm:"not null;type:enum('text','number','boolean');default:text;size:255" json:"type"`
	Description string    `gorm:"type:text" json:"description"`
	IsActive    bool      `gorm:"not null;default:true;index" json:"is_active"`
	CreatedAt   time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt   time.Time `gorm:"not null" json:"updated_at"`
}

func (t *Setting) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == "" {
		t.ID = GenerateID(7)
	}

	return
}
