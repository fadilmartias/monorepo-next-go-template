package models

type PaymentGateway struct {
	BaseModel
	Order          int     `gorm:"not null;index" json:"order"`
	Name           string  `gorm:"not null;size:100" json:"name"`
	Metadata       *string `gorm:"type:text" json:"metadata"`
	Description    *string `gorm:"type:text" json:"description"`
	IsActive       bool    `gorm:"default:true;not null;index" json:"is_active"`
	PaymentMethods []PaymentMethod
}
