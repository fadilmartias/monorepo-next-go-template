package models

import (
	"time"
)

type Banner struct {
	BaseModel
	TenantID   string    `gorm:"not null;size:7" json:"tenant_id"`
	Title      string    `gorm:"not null;size:255" json:"title"`
	Content    *string   `gorm:"type:text" json:"content"`
	Img        *string   `json:"img"`
	Link       *string   `json:"link"`
	IsActive   bool      `gorm:"default:true;not null;index:idx_active_valid" json:"is_active"`
	ValidFrom  time.Time `gorm:"not null;index:idx_active_valid" json:"valid_from"`
	ValidUntil time.Time `gorm:"not null;index:idx_active_valid" json:"valid_until"`
	Order      int       `gorm:"not null;index:idx_active_valid" json:"order"`
}
