package models

import (
	"time"

	"gorm.io/gorm"
)

type Banner struct {
	ID         string    `gorm:"primarykey;size:7" json:"id"`
	TenantID   string    `gorm:"not null;size:7" json:"tenant_id"`
	Title      string    `gorm:"not null;size:100" json:"title"`
	Content    *string   `gorm:"type:text" json:"content"`
	Img        *string   `json:"img"`
	Link       *string   `json:"link"`
	IsActive   bool      `gorm:"default:true;not null;index:idx_active_valid" json:"is_active"`
	ValidFrom  time.Time `gorm:"not null;index:idx_active_valid" json:"valid_from"`
	ValidUntil time.Time `gorm:"not null;index:idx_active_valid" json:"valid_until"`
	Order      int       `gorm:"not null;index:idx_active_valid" json:"order"`
	CreatedAt  time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (t *Banner) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == "" {
		t.ID = GenerateID(7)
	}

	return

}
