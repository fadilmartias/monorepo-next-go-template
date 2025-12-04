package dto

type PaymentMethodDTO struct {
	ID                string  `json:"id"`
	PaymentGatewayID  string  `json:"payment_gateway_id"`
	CategoryID        string  `json:"category_id"`
	Order             int     `json:"order"`
	Code              string  `json:"code"`
	InvoiceCode       string  `json:"invoice_code"`
	Name              string  `json:"name"`
	FeeFixed          float64 `json:"fee_fixed"`
	FeePercent        float64 `json:"fee_percent"`
	PPN               float64 `json:"ppn"`
	MinAmount         float64 `json:"min_amount"`
	MaxAmount         float64 `json:"max_amount"`
	Img               *string `json:"img"`
	Desc              *string `json:"desc"`
	IsActive          bool    `json:"is_active"`
	IsReadyProduction bool    `json:"is_ready_production"`
	IsInternational   bool    `json:"is_international"`
	CategoryName      string  `json:"category_name"`
}
