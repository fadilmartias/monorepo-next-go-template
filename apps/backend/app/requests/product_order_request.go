package requests

type CheckoutInput struct {
	ProductID        string `json:"product_id" validate:"required,min=1"`
	ProductVariantID string `json:"product_variant_id" validate:"required,min=1"`
	// promotions sekarang array of object
	Promotions     []PromotionInput `json:"promotions,omitempty"`
	PreorderRuleID string           `json:"preorder_rule_id,omitempty"`
	GameUsername   string           `json:"game_username,omitempty"`
	Qty            int              `json:"qty" validate:"required,min=1"`
	WaNumber       string           `json:"wa_number" validate:"required,min=10"`
	PaymentMethod  struct {
		ID              string `json:"id" validate:"required,min=1"`
		CardNumber      string `json:"card_number,omitempty"`
		CardExpiryMonth string `json:"card_expiry_month,omitempty"`
		CardExpiryYear  string `json:"card_expiry_year,omitempty"`
		CVV             string `json:"cvv,omitempty"`
		MobileNumber    string `json:"mobile_number,omitempty"`
	} `json:"payment_method,omitempty"`
	AdditionalFields map[string]string `json:"additional_fields,omitempty"`
}

type PromotionInput struct {
	ID          string `json:"id" validate:"required,min=1"`
	Type        string `json:"type" validate:"required"`
	Code        string `json:"code,omitempty"`
	IsRedeemed  bool   `json:"is_redeemed,omitempty"`
	Description string `json:"description,omitempty"`
}
