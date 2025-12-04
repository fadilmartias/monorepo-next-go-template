package responses

type PaymentGatewayResponse struct {
	ID        *int    `json:"id,omitempty"`
	CreatedAt *string `json:"created_at,omitempty"`
	UpdatedAt *string `json:"updated_at,omitempty"`
}
