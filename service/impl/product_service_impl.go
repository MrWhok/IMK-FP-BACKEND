package impl

import (
	"context"
	"github.com/MrWhok/IMK-FP-BACKEND/common"
	"github.com/MrWhok/IMK-FP-BACKEND/configuration"
	"github.com/MrWhok/IMK-FP-BACKEND/entity"
	"github.com/MrWhok/IMK-FP-BACKEND/exception"
	"github.com/MrWhok/IMK-FP-BACKEND/model"
	"github.com/MrWhok/IMK-FP-BACKEND/repository"
	"github.com/MrWhok/IMK-FP-BACKEND/service"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
)

func NewProductServiceImpl(productRepository *repository.ProductRepository, cache *redis.Client) service.ProductService {
	return &productServiceImpl{ProductRepository: *productRepository, Cache: cache}
}

type productServiceImpl struct {
	repository.ProductRepository
	Cache *redis.Client
}

func (service *productServiceImpl) Create(ctx context.Context, productModel model.ProductCreateOrUpdateModel, imagePath string) model.ProductCreateOrUpdateModel {
	common.Validate(productModel)
	product := entity.Product{
		Name:     productModel.Name,
		Price:    productModel.Price,
		Quantity: productModel.Quantity,
		ImagePath:    imagePath,
	}
	service.ProductRepository.Insert(ctx, product)
	return productModel
}

func (service *productServiceImpl) Update(ctx context.Context, productModel model.ProductCreateOrUpdateModel, id string) model.ProductCreateOrUpdateModel {
	common.Validate(productModel)
	product := entity.Product{
		Id:       uuid.MustParse(id),
		Name:     productModel.Name,
		Price:    productModel.Price,
		Quantity: productModel.Quantity,
	}
	service.ProductRepository.Update(ctx, product)
	return productModel
}

func (service *productServiceImpl) Delete(ctx context.Context, id string) {
	product, err := service.ProductRepository.FindById(ctx, id)
	if err != nil {
		panic(exception.NotFoundError{
			Message: err.Error(),
		})
	}
	service.ProductRepository.Delete(ctx, product)
}

func (service *productServiceImpl) FindById(ctx context.Context, id string) model.ProductModel {
	productCache := configuration.SetCache[entity.Product](service.Cache, ctx, "product", id, service.ProductRepository.FindById)
	return model.ProductModel{
		Id:       productCache.Id.String(),
		Name:     productCache.Name,
		Price:    productCache.Price,
		Quantity: productCache.Quantity,
		ImagePath:    productCache.ImagePath,
	}
}

func (service *productServiceImpl) FindAll(ctx context.Context) (responses []model.ProductModel) {
	products := service.ProductRepository.FindAl(ctx)
	for _, product := range products {
		responses = append(responses, model.ProductModel{
			Id:       product.Id.String(),
			Name:     product.Name,
			Price:    product.Price,
			Quantity: product.Quantity,
			ImagePath:    product.ImagePath,
		})
	}
	if len(products) == 0 {
		return []model.ProductModel{}
	}
	return responses
}
