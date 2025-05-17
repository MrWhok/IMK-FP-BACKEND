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
		panic(exception.NotFoundError{
			Message: err.Error(),
		})
	}
	if product.Quantity < request.Quantity {
		panic(exception.BadRequestError{
			Message: "Not enough product quantity",
		})
	}

	entityCart := entity.Cart{
		Username:  username,
		ProductID: request.ProductID,
		Quantity:  request.Quantity,
	}

	result := s.cartRepo.AddToCart(ctx, entityCart)

	return model.AddToCartResponse{
		ID:          result.ID,
		ProductID:   product.Id.String(),
		ProductName: product.Name,
		Quantity:    result.Quantity,
	}
}

func (s *cartServiceImpl) GetMyCart(ctx context.Context, username string) model.CartItemFinalResponse {
	carts, _ := s.cartRepo.FindByUsername(ctx, username)

	var items []model.CartItemResponse
	for _, item := range carts {
		items = append(items, model.CartItemResponse{
			ProductID: item.ProductID,
			Name:      item.Product.Name,
			Price:     item.Product.Price,
			Quantity:  item.Quantity,
			ImagePath: item.Product.ImagePath,
		})
	}

	return model.CartItemFinalResponse{
		Username: username,
		Items:    items,
	}
}
