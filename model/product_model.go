package model

import (
	"mime/multipart"
)

type ProductModel struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Price       int64  `json:"price"`
	Quantity    int32  `json:"quantity"`
	Category    string `json:"category"`
	Description string `json:"description"`
	ImagePath   string `json:"image_path"`
	Owner       string `json:"owner"`
	OwnerPhone  string `json:"owner_phone"`
}

type ProductCreateModel struct {
	Name        string                `form:"name" validate:"required"`
	Price       int64                 `form:"price" validate:"required"`
	Quantity    int32                 `form:"quantity" validate:"required"`
	Category    string                `form:"category" validate:"required"`
	Description string                `form:"description"`
	Image       *multipart.FileHeader `form:"image" validate:"required"`
}

type ProductUpdateModel struct {
	Name        string                `form:"name"`
	Price       int64                 `form:"price"`
	Quantity    int32                 `form:"quantity"`
	Category    string                `form:"category"`
	Description string                `form:"description"`
	Image       *multipart.FileHeader `form:"image"`
}
