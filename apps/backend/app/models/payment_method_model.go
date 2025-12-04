package models

import (
	"time"

	"gorm.io/gorm"
)

type PaymentMethod struct {
	ID                string  `gorm:"primaryKey;size:7" json:"id"`
	PaymentGatewayID  string  `gorm:"not null;size:7;index" json:"payment_gateway_id"`
	CategoryID        string  `gorm:"not null;size:7;index" json:"category_id"`
	Order             int     `gorm:"not null;index" json:"order"`
	Type              string  `gorm:"not null;" json:"type"`
	Code              string  `gorm:"size:100" json:"code"`
	InvoiceCode       string  `gorm:"size:100" json:"invoice_code"`
	Name              string  `gorm:"not null;size:100" json:"name"`
	FeeFixed          float64 `gorm:"default:0;not null" json:"fee_fixed"`
	FeePercent        float64 `gorm:"default:0;not null" json:"fee_percent"`
	PPN               float64 `gorm:"default:0;not null" json:"ppn"`
	MinAmount         float64 `gorm:"default:0;not null" json:"min_amount"`
	MaxAmount         float64 `gorm:"default:0;not null" json:"max_amount"`
	Img               *string `gorm:"size:255" json:"img"`
	ExpirySeconds     int     `gorm:"default:60;not null" json:"expiry_seconds"`
	Desc              *string `gorm:"type:text" json:"desc"`
	IsActive          bool    `gorm:"default:true;not null;index" json:"is_active"`
	IsReadyProduction bool    `gorm:"default:false;not null;index" json:"is_ready_production"`
	IsInternational   bool    `gorm:"default:false;not null;index" json:"is_international"`
	AdditionalFields  *string `gorm:"type:text" json:"additional_fields"` // could be JSON
	Category          *Category
	CreatedAt         time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt         time.Time `gorm:"not null" json:"updated_at"`
	PaymentGateway    *PaymentGateway
}

func (t *PaymentMethod) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == "" {
		t.ID = GenerateID(7)
	}

	return
}
