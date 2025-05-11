package service

import (
	"context"
	"github.com/MrWhok/IMK-FP-BACKEND/model"
)

type TransactionDetailService interface {
	FindById(ctx context.Context, id string) model.TransactionDetailModel
}
