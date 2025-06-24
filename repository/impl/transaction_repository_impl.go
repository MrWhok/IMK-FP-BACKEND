package impl

import (
	"context"
	"errors"

	"github.com/MrWhok/IMK-FP-BACKEND/entity"
	"github.com/MrWhok/IMK-FP-BACKEND/exception"
	"github.com/MrWhok/IMK-FP-BACKEND/repository"
	"gorm.io/gorm"
)

func NewTransactionRepositoryImpl(DB *gorm.DB) repository.TransactionRepository {
	return &transactionRepositoryImpl{DB: DB}
}

type transactionRepositoryImpl struct {
	*gorm.DB
}

func (transactionRepository *transactionRepositoryImpl) Insert(ctx context.Context, transaction entity.Transaction) entity.Transaction {
	err := transactionRepository.DB.WithContext(ctx).Create(&transaction).Error
	exception.PanicLogging(err)
	return transaction
}

func (transactionRepository *transactionRepositoryImpl) Delete(ctx context.Context, transaction entity.Transaction) {
	transactionRepository.DB.WithContext(ctx).Delete(&transaction)
}

func (transactionRepository *transactionRepositoryImpl) FindById(ctx context.Context, id string) (entity.Transaction, error) {
	var transaction entity.Transaction
	result := transactionRepository.DB.WithContext(ctx).
		Preload("TransactionDetails").
		Preload("TransactionDetails.Product").
		Preload("TransactionDetails.Product.Owner"). // ✅ Tambahan penting
		Where("transaction_id = ?", id).
		First(&transaction)

	if result.RowsAffected == 0 {
		return entity.Transaction{}, errors.New("transaction Not Found")
	}
	return transaction, nil
}

func (transactionRepository *transactionRepositoryImpl) FindAll(ctx context.Context) []entity.Transaction {
	var transactions []entity.Transaction
	transactionRepository.DB.WithContext(ctx).
		Preload("TransactionDetails").
		Preload("TransactionDetails.Product").
		Preload("TransactionDetails.Product.Owner"). // ✅ Tambahan penting
		Find(&transactions)
	return transactions
}

func (transactionRepository *transactionRepositoryImpl) FindByUsername(ctx context.Context, username string) []entity.Transaction {
	var transactions []entity.Transaction
	transactionRepository.DB.WithContext(ctx).
		Table("tb_transaction").
		Select("DISTINCT tb_transaction.transaction_id, tb_transaction.total_price, tb_transaction.user_id, tb_transaction.status").
		Joins("JOIN tb_transaction_detail ON tb_transaction_detail.transaction_id = tb_transaction.transaction_id").
		Joins("JOIN tb_product ON tb_product.product_id = tb_transaction_detail.product_id").
		Joins("JOIN tb_user ON tb_user.username = tb_product.user_id").
		Where("tb_user.username = ?", username).
		Find(&transactions)

	for i := range transactions {
		transactionRepository.DB.WithContext(ctx).
			Preload("TransactionDetails").
			Preload("TransactionDetails.Product").
			Preload("TransactionDetails.Product.Owner"). // ✅ Tambahan penting
			First(&transactions[i], transactions[i].Id)
	}

	return transactions
}

func (transactionRepository *transactionRepositoryImpl) FindByBuyerUsername(ctx context.Context, username string) []entity.Transaction {
	var transactions []entity.Transaction
	transactionRepository.DB.WithContext(ctx).
		Preload("TransactionDetails").
		Preload("TransactionDetails.Product").
		Preload("TransactionDetails.Product.Owner"). // ✅ Tambahan penting
		Where("user_id = ?", username).
		Find(&transactions)

	return transactions
}

func (transactionRepository *transactionRepositoryImpl) UpdateStatus(ctx context.Context, id string, status string) error {
	result := transactionRepository.DB.WithContext(ctx).
		Model(&entity.Transaction{}).
		Where("transaction_id = ?", id).
		Update("status", status)

	if result.RowsAffected == 0 {
		return errors.New("transaction Not Found")
	}
	if result.Error != nil {
		return result.Error
	}
	return nil
}
