package impl

import (
	"context"

	"github.com/MrWhok/IMK-FP-BACKEND/entity"
	"github.com/MrWhok/IMK-FP-BACKEND/exception"
	"github.com/MrWhok/IMK-FP-BACKEND/model"
	"github.com/MrWhok/IMK-FP-BACKEND/repository"
	"github.com/MrWhok/IMK-FP-BACKEND/service"
)

func NewCartServiceImpl(cartRepo repository.CartRepository, productRepo repository.ProductRepository) service.CartService {
	return &cartServiceImpl{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

type cartServiceImpl struct {
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
}

func (s *cartServiceImpl) AddToCart(ctx context.Context, username string, request model.AddToCartRequest) model.AddToCartResponse {
	product, err := s.productRepo.FindById(ctx, request.ProductID)
	if err != nil {
		panic(exception.NotFoundError{Message: err.Error()})
	}
	if product.Quantity < request.Quantity {
		panic(exception.BadRequestError{Message: "Not enough product quantity"})
	}

	cart, _ := s.cartRepo.FindOrCreateCartByUsername(ctx, username)

	existingItem, err := s.cartRepo.FindItemByUsernameAndProductID(ctx, username, request.ProductID)
	if err == nil {
		// Update quantity
		existingItem.Quantity += request.Quantity
		s.cartRepo.UpdateItem(ctx, existingItem)

		return model.AddToCartResponse{
			ItemID:      existingItem.ID,
			CartID:      cart.ID,
			ProductID:   product.Id.String(),
			ProductName: product.Name,
			Quantity:    existingItem.Quantity,
		}
	}

	// Insert new item
	newItem := entity.CartItem{
		CartID:    cart.ID,
		ProductID: request.ProductID,
		Quantity:  request.Quantity,
	}
	result := s.cartRepo.AddToCartItem(ctx, newItem)

	return model.AddToCartResponse{
		ItemID:      result.ID,
		CartID:      cart.ID,
		ProductID:   product.Id.String(),
		ProductName: product.Name,
		Quantity:    result.Quantity,
	}
}

func (s *cartServiceImpl) GetMyCart(ctx context.Context, username string) model.CartItemFinalResponse {
	items, _ := s.cartRepo.FindCartItemsByUsername(ctx, username)

	var cartItems []model.CartItemResponse
	for _, item := range items {
		cartItems = append(cartItems, model.CartItemResponse{
			ProductID: item.ProductID,
			Name:      item.Product.Name,
			Price:     item.Product.Price,
			Quantity:  item.Quantity,
			ImagePath: item.Product.ImagePath,
		})
	}

	return model.CartItemFinalResponse{
		Username: username,
		Items:    cartItems,
	}
}

func (s *cartServiceImpl) UpdateCartItem(ctx context.Context, username string, productID string, req model.UpdateCartRequest) {
	cartItem, err := s.cartRepo.FindItemByUsernameAndProductID(ctx, username, productID)
	if err != nil {
		panic(exception.NotFoundError{Message: "Cart item not found"})
	}

	cartItem.Quantity = req.Quantity
	s.cartRepo.UpdateItem(ctx, cartItem)
}

func (s *cartServiceImpl) DeleteCartItem(ctx context.Context, username string, productID string) {
	s.cartRepo.DeleteItem(ctx, username, productID)
}

func (s *cartServiceImpl) SubstractFromCart(ctx context.Context, username string, productID string) {
	cartItem, err := s.cartRepo.FindItemByUsernameAndProductID(ctx, username, productID)
	if err != nil {
		panic(exception.NotFoundError{Message: "Cart item not found"})
	}

	if cartItem.Quantity <= 1 {
		s.cartRepo.DeleteItem(ctx, username, productID)
		return
	}

	cartItem.Quantity -= 1
	s.cartRepo.UpdateItem(ctx, cartItem)
}
