package model

type AddToCartRequest struct {
	ProductID string `json:"product_id" validate:"required,uuid4"`
	Quantity  int32  `json:"quantity" validate:"required,min=1"`
}

type AddToCartResponse struct {
	ID          uint   `json:"id"`
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	Quantity    int32  `json:"quantity"`
}

type CartItemResponse struct {
	ProductID string `json:"product_id"`
	Name      string `json:"name"`
	Price     int64  `json:"price"`
	Quantity  int32  `json:"quantity"`
	ImagePath string `json:"image_path"`
}

type CartItemFinalResponse struct {
	Username string             `json:"username"`
	Items    []CartItemResponse `json:"items"`
}
