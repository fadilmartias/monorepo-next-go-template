package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID             string     `gorm:"primaryKey;size:7" json:"id"`
	ParentID       *string    `gorm:"size:7;index" json:"parent_id"`
	Parent         *Category  `gorm:"foreignKey:ParentID"`
	Children       []Category `gorm:"foreignKey:ParentID"`
	Order          int        `gorm:"not null;index" json:"order"`
	Name           string     `gorm:"not null;size:100" json:"name"`
	Desc           *string    `gorm:"type:text" json:"desc"`
	Slug           string     `gorm:"not null;size:100;uniqueIndex" json:"slug"`
	IsActive       bool       `gorm:"default:true;not null;index" json:"is_active"`
	Type           string     `gorm:"type:enum('product','transaction','payment-method','article','product-subcategory','product-parent');not null;index" json:"type"`
	CreatedAt      time.Time  `gorm:"not null" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"not null" json:"updated_at"`
	PaymentMethods []PaymentMethod
}

func (t *Category) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == "" {
		t.ID = GenerateID(7)
	}

	return
}
