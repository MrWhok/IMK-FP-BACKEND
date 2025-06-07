package model

import (
	"mime/multipart"
)

type GiftModel struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	PointPrice int64  `json:"point_price"`
	Quantity   int32  `json:"quantity"`
	ImagePath  string `json:"image_path"`
}

type GiftCreateModel struct {
	Name       string                `form:"name" validate:"required"`
	PointPrice int64                 `form:"point_price" validate:"required"`
	Quantity   int32                 `form:"quantity" validate:"required"`
	Image      *multipart.FileHeader `form:"image" validate:"required"`
}

type GiftUpdateModel struct {
	Name       *string               `form:"name"`
	PointPrice *int64                `form:"point_price"`
	Quantity   *int32                `form:"quantity"`
	Image      *multipart.FileHeader `form:"image"`
}
