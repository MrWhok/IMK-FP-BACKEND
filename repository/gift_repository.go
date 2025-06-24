package repository

import (
	"context"

	"github.com/MrWhok/IMK-FP-BACKEND/entity"
)

type GiftRepository interface {
	Insert(ctx context.Context, gift entity.Gift) entity.Gift
	Update(ctx context.Context, gift entity.Gift) entity.Gift
	Delete(ctx context.Context, gift entity.Gift)
	FindById(ctx context.Context, id string) (entity.Gift, error)
	FindAll(ctx context.Context) []entity.Gift
	ExchangeGift(ctx context.Context, giftId string, username string) error
}
