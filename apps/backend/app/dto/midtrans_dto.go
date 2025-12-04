package dto

import "github.com/fadilmartias/dilz_code/apps/backend/app/models"

type CreateTransactionParams struct {
	ID                  string                   `json:"id"`
	Qty                 int                      `json:"qty"`
	PaymentMethod       models.PaymentMethod     `json:"payment_method"`
	BasePrice           float64                  `json:"base_price"`
	PriceAfterMargin    float64                  `json:"price_after_margin"`
	TotalPrice          float64                  `json:"total_price"`
	ItemPrice           float64                  `json:"item_price"`
	ItemPricePlusMargin float64                  `json:"item_price_plus_margin"`
	TotalFee            float64                  `json:"total_fee"`
	PGFee               float64                  `json:"pg_fee"`
	Profit              float64                  `json:"profit"`
	TotalDiscount       *float64                 `json:"total_discount,omitempty"`
	OrderId             string                   `json:"order_id"`
	RefId               string                   `json:"ref_id"`
	NoWA                string                   `json:"no_wa"`
	CustomerDetails     *MidtransCustomerDetails `json:"customer_details,omitempty"`
}

type MidtransCustomerDetails struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}
