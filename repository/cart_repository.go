package repository

import (
	"context"

	"github.com/MrWhok/IMK-FP-BACKEND/entity"
)

type CartRepository interface {
	AddToCartItem(ctx context.Context, item entity.CartItem) entity.CartItem
	FindOrCreateCartByUsername(ctx context.Context, username string) (entity.Cart, error)
	FindCartItemsByUsername(ctx context.Context, username string) ([]entity.CartItem, error)
	FindItemByUsernameAndProductID(ctx context.Context, username, productID string) (entity.CartItem, error)
	UpdateItem(ctx context.Context, item entity.CartItem)
	DeleteItem(ctx context.Context, username, productID string)
}
