package service

import (
	"context"

	"github.com/MrWhok/IMK-FP-BACKEND/model"
)

type GiftService interface {
	Create(ctx context.Context, model model.GiftCreateModel, imagePath string) model.GiftCreateModel
	Update(ctx context.Context, giftModel model.GiftUpdateModel, id string) model.GiftModel
	Delete(ctx context.Context, id string) error
	FindById(ctx context.Context, id string) model.GiftModel
	FindAll(ctx context.Context) []model.GiftModel
	ExchangeGift(ctx context.Context, giftId string, username string) error
}
