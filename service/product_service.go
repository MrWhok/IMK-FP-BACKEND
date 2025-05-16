package service

import (
	"context"

	"github.com/MrWhok/IMK-FP-BACKEND/model"
)

type ProductService interface {
	Create(ctx context.Context, model model.ProductCreateModel, imagePath string) model.ProductCreateModel
	Update(ctx context.Context, productModel model.ProductUpdateModel, id string) model.ProductModel
	Delete(ctx context.Context, id string) error
	FindById(ctx context.Context, id string) model.ProductModel
	FindAll(ctx context.Context) []model.ProductModel
}
