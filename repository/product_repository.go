package repository

import (
	"context"

	"github.com/MrWhok/IMK-FP-BACKEND/entity"
)

type ProductRepository interface {
	Insert(ctx context.Context, product entity.Product) entity.Product
	Update(ctx context.Context, product entity.Product) entity.Product
	Delete(ctx context.Context, product entity.Product)
	FindById(ctx context.Context, id string) (entity.Product, error)
	FindAl(ctx context.Context) []entity.Product
	FindByUsername(ctx context.Context, username string) []entity.Product
}
