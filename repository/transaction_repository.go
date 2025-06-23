package repository

import (
	"context"

	"github.com/MrWhok/IMK-FP-BACKEND/entity"
)

type TransactionRepository interface {
	Insert(ctx context.Context, transaction entity.Transaction) entity.Transaction
	Delete(ctx context.Context, transaction entity.Transaction)
	FindById(ctx context.Context, id string) (entity.Transaction, error)
	FindAll(ctx context.Context) []entity.Transaction
	FindByUsername(ctx context.Context, username string) []entity.Transaction
	FindByBuyerUsername(ctx context.Context, username string) []entity.Transaction
}
