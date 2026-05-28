package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        string    `gorm:"type:char(36);primaryKey" json:"id"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BaseModelWithoutUpdatedAt struct {
	ID        string    `gorm:"type:char(36);primaryKey" json:"id"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
}

type BaseModelWithDeletedAt struct {
	ID        string         `gorm:"type:char(36);primaryKey" json:"id"`
	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func insertIDIfMissing(id *string) error {
	if id == nil {
		return nil
	}

	if *id != "" {
		return nil
	}

	uuidV7, err := uuid.NewV7()
	if err != nil {
		return err
	}

	*id = uuidV7.String()
	return nil
}

// BeforeCreate hook untuk generate UUID v7 otomatis
func (b *BaseModel) BeforeCreate(_ *gorm.DB) error {
	return insertIDIfMissing(&b.ID)
}

// BeforeCreate hook untuk generate UUID v7 otomatis
func (b *BaseModelWithoutUpdatedAt) BeforeCreate(_ *gorm.DB) error {
	return insertIDIfMissing(&b.ID)
}

// BeforeCreate hook untuk generate UUID v7 otomatis
func (b *BaseModelWithDeletedAt) BeforeCreate(_ *gorm.DB) error {
	return insertIDIfMissing(&b.ID)
}
