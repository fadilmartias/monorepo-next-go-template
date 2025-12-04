package controllers_v1

import (
	"github.com/fadilmartias/dilz_code/apps/backend/app/services"
	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"

	"github.com/gofiber/fiber/v2"
)

type PaymentMethodController struct {
	BaseController
	PaymentMethodService *services.PaymentMethodService
}

func NewPaymentMethodController(paymentMethodService *services.PaymentMethodService) *PaymentMethodController {
	return &PaymentMethodController{PaymentMethodService: paymentMethodService}
}

type PaymentMethodDTO struct {
	ID                string  `gorm:"primaryKey;size:7" json:"id"`
	PaymentGatewayID  string  `gorm:"not null;size:7;index" json:"payment_gateway_id"`
	CategoryID        string  `gorm:"not null;size:7;index" json:"category_id"`
	Order             int     `gorm:"not null" json:"order"`
	Code              string  `gorm:"not null;unique;size:50" json:"code"`
	InvoiceCode       string  `gorm:"not null;unique;size:50" json:"invoice_code"`
	Name              string  `gorm:"not null;size:50" json:"name"`
	FeeFixed          float64 `gorm:"not null" json:"fee_fixed"`
	FeePercent        float64 `gorm:"not null" json:"fee_percent"`
	PPN               float64 `gorm:"not null" json:"ppn"`
	MinAmount         float64 `gorm:"not null" json:"min_amount"`
	MaxAmount         float64 `gorm:"not null" json:"max_amount"`
	Img               *string `gorm:"size:50" json:"img"`
	Desc              *string `gorm:"size:255" json:"desc"`
	IsActive          bool    `gorm:"not null" json:"is_active"`
	IsReadyProduction bool    `gorm:"not null" json:"is_ready_production"`
	IsInternational   bool    `gorm:"not null" json:"is_international"`
	CategoryName      string  `json:"category_name"`
}

func (ctrl *PaymentMethodController) Index(c *fiber.Ctx) error {
	isCache := c.QueryBool("cache", false)
	cacheTtl := c.QueryInt("cache_ttl", 60*60*24)
	cacheKey := c.Query("cache_key", "active-payment-methods")
	data, err := ctrl.PaymentMethodService.GetActivePaymentMethods(c, isCache, cacheTtl, cacheKey)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Message: "Gagal mendapatkan data payment method",
			Details: err.Error(),
		})
	}

	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Berhasil mendapatkan data payment method",
		Data:    data,
	})
}
