package service

import (
	"context"

	"github.com/MrWhok/IMK-FP-BACKEND/model"
)

type CartService interface {
	AddToCart(ctx context.Context, username string, request model.AddToCartRequest) model.AddToCartResponse
	GetMyCart(ctx context.Context, username string) model.CartItemFinalResponse
}
