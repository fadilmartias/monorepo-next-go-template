package requests

type ProductRatingStoreRequest struct {
	ProductOrderID string `json:"product_order_id" validate:"required,max=255"`
	ProductID      string `json:"product_id" validate:"required,max=255"`
	UserID         string `json:"user_id,omitempty" validate:"max=7"`
	Name           string `json:"name" validate:"required,max=255"`
	Rating         int    `json:"rating" validate:"required"`
	Comment        string `json:"comment,omitempty" validate:"max=255"`
}
