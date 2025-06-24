package impl

import (
	"context"
	"errors"

	"github.com/MrWhok/IMK-FP-BACKEND/entity"
	"github.com/MrWhok/IMK-FP-BACKEND/repository"
	"gorm.io/gorm"
)

func NewTransactionDetailRepositoryImpl(DB *gorm.DB) repository.TransactionDetailRepository {
	return &transactionDetailRepositoryImpl{DB: DB}
}

type transactionDetailRepositoryImpl struct {
	*gorm.DB
}

func (transactionDetailRepository *transactionDetailRepositoryImpl) FindById(ctx context.Context, id string) (entity.TransactionDetail, error) {
	var transactionDetail entity.TransactionDetail
	result := transactionDetailRepository.DB.WithContext(ctx).
		Where("transaction_detail_id = ?", id).
		Preload("Product").
		Preload("Product.Owner"). // ✅ tambahkan ini
		First(&transactionDetail)
	if result.RowsAffected == 0 {
		return entity.TransactionDetail{}, errors.New("transaction Detail Not Found")
	}
	return transactionDetail, nil
}
