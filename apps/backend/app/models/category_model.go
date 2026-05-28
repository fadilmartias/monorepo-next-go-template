package models

type Category struct {
	BaseModel
	ParentID       *string      `gorm:"size:7;index" json:"parent_id"`
	Parent         *Category    `gorm:"foreignKey:ParentID"`
	Children       []Category   `gorm:"foreignKey:ParentID"`
	Order          int          `gorm:"not null;index" json:"order"`
	Name           string       `gorm:"not null;size:100" json:"name"`
	Desc           *string      `gorm:"type:text" json:"desc"`
	Slug           string       `gorm:"not null;size:100;uniqueIndex" json:"slug"`
	IsActive       bool         `gorm:"default:true;not null;index" json:"is_active"`
	Type           CategoryType `gorm:"size:100;not null;index" json:"type"`
	PaymentMethods []PaymentMethod
}
