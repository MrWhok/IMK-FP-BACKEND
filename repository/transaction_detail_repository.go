package repository

import (
	"context"
	"github.com/MrWhok/IMK-FP-BACKEND/entity"
)

type TransactionDetailRepository interface {
	FindById(ctx context.Context, id string) (entity.TransactionDetail, error)
}
