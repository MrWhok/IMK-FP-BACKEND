package impl

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/MrWhok/IMK-FP-BACKEND/common"
	"github.com/MrWhok/IMK-FP-BACKEND/entity"
	"github.com/MrWhok/IMK-FP-BACKEND/exception"
	"github.com/MrWhok/IMK-FP-BACKEND/model"
	"github.com/MrWhok/IMK-FP-BACKEND/repository"
	"github.com/MrWhok/IMK-FP-BACKEND/service"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
)

func NewGiftServiceImpl(giftRepository *repository.GiftRepository, cache *redis.Client) service.GiftService {
	return &giftServiceImpl{GiftRepository: *giftRepository, Cache: cache}
}

type giftServiceImpl struct {
	repository.GiftRepository
	Cache *redis.Client
}

func (service *giftServiceImpl) Create(ctx context.Context, giftModel model.GiftCreateModel, imagePath string) model.GiftCreateModel {
	common.Validate(giftModel)

	gift := entity.Gift{
		Name:       giftModel.Name,
		PointPrice: giftModel.PointPrice,
		Quantity:   giftModel.Quantity,
		ImagePath:  imagePath,
	}
	service.GiftRepository.Insert(ctx, gift)
	return giftModel
}

func (service *giftServiceImpl) Update(ctx context.Context, giftModel model.GiftUpdateModel, id string) model.GiftModel {
	// Validasi input
	common.Validate(giftModel)

	// Ambil produk lama dari database
	existingGift, err := service.GiftRepository.FindById(ctx, id)
	if err != nil {
		panic(exception.NotFoundError{Message: err.Error()})
	}

	// Update field dasar
	if giftModel.Name != nil {
		existingGift.Name = *giftModel.Name
	}
	if giftModel.PointPrice != nil {
		existingGift.PointPrice = *giftModel.PointPrice
	}
	if giftModel.Quantity != nil {
		existingGift.Quantity = *giftModel.Quantity
	}

	// Jika ada image baru
	if giftModel.Image != nil {
		fmt.Println("New image uploaded, deleting old and saving new image...")

		// Hapus gambar lama
		service.deleteGiftImage(existingGift.ImagePath)

		// Simpan gambar baru
		imageID := uuid.New().String()
		imageName := fmt.Sprintf("%s.png", imageID)
		imagePath := fmt.Sprintf("./media/gifts/%s", imageName)

		src, err := giftModel.Image.Open()
		if err != nil {
			panic(exception.InternalServerError{Message: "Failed to open new image"})
		}
		defer src.Close()

		dst, err := os.Create(imagePath)
		if err != nil {
			panic(exception.InternalServerError{Message: "Failed to create image file"})
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			panic(exception.InternalServerError{Message: "Failed to save image"})
		}

		existingGift.ImagePath = imagePath
		fmt.Println("New image saved:", imagePath)
	} else {
		fmt.Println("No new image provided")
	}

	// Update ke database
	service.GiftRepository.Update(ctx, existingGift)

	// Hapus cache Redis agar tidak ambil data lama
	service.Cache.Del(ctx, "gift:"+id)

	// Return response
	return model.GiftModel{
		Id:         existingGift.Id.String(),
		Name:       existingGift.Name,
		PointPrice: existingGift.PointPrice,
		Quantity:   existingGift.Quantity,
		ImagePath:  existingGift.ImagePath,
	}
}

func (service *giftServiceImpl) Delete(ctx context.Context, id string) error {
	gift, err := service.GiftRepository.FindById(ctx, id)
	if err != nil {
		panic(exception.NotFoundError{
			Message: err.Error(),
		})
	}

	service.deleteGiftImage(gift.ImagePath)

	service.GiftRepository.Delete(ctx, gift)
	return nil
}

func (service *giftServiceImpl) deleteGiftImage(imagePath string) {
	if imagePath != "" {
		if err := os.Remove(imagePath); err != nil && !os.IsNotExist(err) {
			panic(exception.InternalServerError{Message: "Failed to delete image"})
		}
	}
}

func (service *giftServiceImpl) FindById(ctx context.Context, id string) model.GiftModel {
	gift, _ := service.GiftRepository.FindById(ctx, id)
	return model.GiftModel{
		Id:         gift.Id.String(),
		Name:       gift.Name,
		PointPrice: gift.PointPrice,
		Quantity:   gift.Quantity,
		ImagePath:  gift.ImagePath,
	}
}

func (service *giftServiceImpl) FindAll(ctx context.Context) (responses []model.GiftModel) {
	gifts := service.GiftRepository.FindAll(ctx)
	for _, gift := range gifts {
		responses = append(responses, model.GiftModel{
			Id:         gift.Id.String(),
			Name:       gift.Name,
			PointPrice: gift.PointPrice,
			Quantity:   gift.Quantity,
			ImagePath:  gift.ImagePath,
		})
	}
	if len(gifts) == 0 {
		return []model.GiftModel{}
	}
	return responses
}

func (service *giftServiceImpl) ExchangeGift(ctx context.Context, giftId string, username string) error {
	err := service.GiftRepository.ExchangeGift(ctx, giftId, username)
	if err != nil {
		panic(exception.InternalServerError{Message: err.Error()})
	}

	return nil
}
