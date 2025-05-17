package repository

import (
	"context"

	"github.com/MrWhok/IMK-FP-BACKEND/entity"
)

type CartRepository interface {
	AddToCart(ctx context.Context, cart entity.Cart) entity.Cart
	FindByUsername(ctx context.Context, username string) ([]entity.Cart, error)
}
