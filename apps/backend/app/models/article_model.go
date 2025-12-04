package models

import (
	"time"

	"gorm.io/gorm"
)

type Article struct {
	ID              string    `gorm:"primarykey;size:7" json:"id,omitempty"`
	TenantID        string    `gorm:"not null;size:7" json:"tenant_id,omitempty"`
	CategoryID      string    `gorm:"not null;size:7;index" json:"category_id,omitempty"`
	Category        Category  `gorm:"foreignKey=CategoryID;references:ID" json:"category,omitempty"`
	Title           string    `gorm:"not null;size:255" json:"title,omitempty"`
	Author          string    `gorm:"not null;size:255" json:"author,omitempty"`
	Slug            string    `gorm:"not null;size:255;uniqueIndex" json:"slug,omitempty"`
	Content         string    `gorm:"type:text;not null" json:"content,omitempty"`
	Img             string    `gorm:"size:255" json:"img,omitempty"`
	Views           int       `gorm:"default:0;not null" json:"views,omitempty"`
	Likes           int       `gorm:"default:0;not null" json:"likes,omitempty"`
	MetaDescription *string   `gorm:"type:text" json:"meta_description,omitempty"`
	MetaKeywords    *string   `gorm:"size:255" json:"meta_keywords,omitempty"`
	MetaImg         *string   `gorm:"size:255" json:"meta_img,omitempty"`
	PublishedAt     time.Time `gorm:"not null;index:idx_publish" json:"published_at,omitempty"`
	IsActive        bool      `gorm:"default:true;not null;index:idx_publish" json:"is_active,omitempty"`
	CreatedAt       time.Time `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt       time.Time `json:"updated_at,omitempty"`
}

func (t *Article) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == "" {
		t.ID = GenerateID(7)
	}

	return
}
