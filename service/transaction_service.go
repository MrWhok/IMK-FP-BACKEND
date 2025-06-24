package service

import (
	"context"

	"github.com/MrWhok/IMK-FP-BACKEND/model"
)

type TransactionService interface {
	Create(ctx context.Context, model model.TransactionCreateUpdateModel) model.TransactionCreateUpdateModel
	Delete(ctx context.Context, id string)
	FindById(ctx context.Context, id string) model.TransactionModel
	FindAll(ctx context.Context) []model.TransactionModel
	Checkout(ctx context.Context, username string) []model.TransactionModel
	FindByUsername(ctx context.Context, username string) []model.TransactionModel
	FindByBuyerUsername(ctx context.Context, username string) []model.TransactionModel
	UpdateStatus(ctx context.Context, id string, status string) error
}
