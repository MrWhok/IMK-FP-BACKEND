package model

import "github.com/google/uuid"

type TransactionModel struct {
	Id                 string                   `json:"id"`
	TotalPrice         int64                    `json:"total_price"`
	TransactionDetails []TransactionDetailModel `json:"transaction_details"`
}

type TransactionCreateUpdateModel struct {
	Id                 string                               `json:"id"`
	TotalPrice         int64                                `json:"total_price"`
	TransactionDetails []TransactionDetailCreateUpdateModel `json:"transaction_details"`
}

type TransactionDetailModel struct {
	Id            string `json:"id"`
	SubTotalPrice int64  `json:"sub_total_price" validate:"required"`
	Price         int64  `json:"price" validate:"required"`
	Quantity      int32  `json:"quantity" validate:"required"`
	Product       ProductModel
}

type TransactionDetailCreateUpdateModel struct {
	Id            string    `json:"id"`
	SubTotalPrice int64     `json:"sub_total_price" validate:"required"`
	Price         int64     `json:"price" validate:"required"`
	Quantity      int32     `json:"quantity" validate:"required"`
	ProductId     uuid.UUID `json:"product_id" validate:"required"`
	Product       ProductModel
}

type CheckoutResponse struct {
	TransactionId string `json:"transaction_id"`
	TotalPrice    int64  `json:"total_price"`
}

type TransactionItemInfo struct {
	ProductID string `json:"product_id"`
	Name      string `json:"name"`
	Price     int64  `json:"price"`
	Quantity  int32  `json:"quantity"`
}
