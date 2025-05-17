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

func (r *cartRepositoryImpl) FindOrCreateCartByUsername(ctx context.Context, username string) (entity.Cart, error) {
	var cart entity.Cart
	err := r.WithContext(ctx).Where("username = ?", username).FirstOrCreate(&cart, entity.Cart{Username: username}).Error
	exception.PanicLogging(err)
	return cart, nil
}

func (r *cartRepositoryImpl) AddToCartItem(ctx context.Context, item entity.CartItem) entity.CartItem {
	err := r.WithContext(ctx).Create(&item).Error
	exception.PanicLogging(err)
	return item
}

func (r *cartRepositoryImpl) FindCartItemsByUsername(ctx context.Context, username string) ([]entity.CartItem, error) {
	var items []entity.CartItem
	err := r.WithContext(ctx).
		Joins("JOIN tb_cart ON tb_cart.cart_id = tb_cart_item.cart_id").
		Preload("Product").
		Where("tb_cart.username = ?", username).
		Find(&items).Error
	exception.PanicLogging(err)
	return items, nil
}

func (r *cartRepositoryImpl) FindItemByUsernameAndProductID(ctx context.Context, username, productID string) (entity.CartItem, error) {
	var item entity.CartItem

	err := r.WithContext(ctx).
		Joins("JOIN tb_cart ON tb_cart.cart_id = tb_cart_item.cart_id").
		Where("tb_cart.username = ? AND tb_cart_item.product_id = ?", username, productID).
		First(&item).Error

	return item, err
}

func (r *cartRepositoryImpl) UpdateItem(ctx context.Context, item entity.CartItem) {
	err := r.WithContext(ctx).Save(&item).Error
	exception.PanicLogging(err)
}

func (r *cartRepositoryImpl) DeleteItem(ctx context.Context, username, productID string) {
	// Get the cart ID first
	var cart entity.Cart
	err := r.WithContext(ctx).
		Where("username = ?", username).
		First(&cart).Error
	if err != nil {
		exception.PanicLogging(err)
		return
	}

	// Now delete the cart item by cart_id and product_id
	err = r.WithContext(ctx).
		Where("cart_id = ? AND product_id = ?", cart.ID, productID).
		Delete(&entity.CartItem{}).Error
	exception.PanicLogging(err)
}
