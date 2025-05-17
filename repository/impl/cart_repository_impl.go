package impl

import (
	"context"

	"github.com/MrWhok/IMK-FP-BACKEND/entity"
	"github.com/MrWhok/IMK-FP-BACKEND/exception"
	"github.com/MrWhok/IMK-FP-BACKEND/repository"
	"gorm.io/gorm"
)

func NewCartRepositoryImpl(DB *gorm.DB) repository.CartRepository {
	return &cartRepositoryImpl{DB: DB}
}

type cartRepositoryImpl struct {
	*gorm.DB
}

func (r *cartRepositoryImpl) AddToCart(ctx context.Context, cart entity.Cart) entity.Cart {
	err := r.WithContext(ctx).Create(&cart).Error
	exception.PanicLogging(err)
	return cart
}

func (r *cartRepositoryImpl) FindByUsername(ctx context.Context, username string) ([]entity.Cart, error) {
	var carts []entity.Cart
	err := r.WithContext(ctx).
		Preload("Product").
		Where("username = ?", username).
		Find(&carts).Error

	if err != nil {
		exception.PanicLogging(err)
	}

	return carts, nil
}
