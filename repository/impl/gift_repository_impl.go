package impl

import (
	"context"
	"errors"

	"github.com/MrWhok/IMK-FP-BACKEND/entity"
	"github.com/MrWhok/IMK-FP-BACKEND/exception"
	"github.com/MrWhok/IMK-FP-BACKEND/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func NewGiftRepositoryImpl(DB *gorm.DB) repository.GiftRepository {
	return &giftRepositoryImpl{DB: DB}
}

type giftRepositoryImpl struct {
	*gorm.DB
}

func (repository *giftRepositoryImpl) Insert(ctx context.Context, gift entity.Gift) entity.Gift {
	gift.Id = uuid.New()
	err := repository.DB.WithContext(ctx).Create(&gift).Error
	exception.PanicLogging(err)
	return gift
}

func (repository *giftRepositoryImpl) Update(ctx context.Context, gift entity.Gift) entity.Gift {
	err := repository.DB.WithContext(ctx).Where("gift_id = ?", gift.Id).Updates(&gift).Error
	exception.PanicLogging(err)
	return gift
}

func (repository *giftRepositoryImpl) Delete(ctx context.Context, gift entity.Gift) {
	err := repository.DB.WithContext(ctx).Delete(&gift).Error
	exception.PanicLogging(err)
}

func (repository *giftRepositoryImpl) FindById(ctx context.Context, id string) (entity.Gift, error) {
	var gift entity.Gift
	result := repository.DB.WithContext(ctx).Unscoped().Where("gift_id = ?", id).First(&gift)
	if result.RowsAffected == 0 {
		return entity.Gift{}, errors.New("gift Not Found")
	}
	return gift, nil
}

func (repository *giftRepositoryImpl) FindAll(ctx context.Context) []entity.Gift {
	var gifts []entity.Gift
	repository.DB.WithContext(ctx).Find(&gifts)
	return gifts
}
