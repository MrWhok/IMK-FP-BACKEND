package impl

import (
	"context"
	"fmt"
	"io"
	"os"

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

func (service *productServiceImpl) Create(ctx context.Context, productModel model.ProductCreateModel, imagePath string) model.ProductCreateModel {
	common.Validate(productModel)

	username := ctx.Value("username").(string)

	product := entity.Product{
		Name:      productModel.Name,
		Price:     productModel.Price,
		Quantity:  productModel.Quantity,
		Category:  productModel.Category,
		ImagePath: imagePath,
		Owner:     entity.User{Username: username},
	}
	service.ProductRepository.Insert(ctx, product)
	return productModel
}

func (service *productServiceImpl) Update(ctx context.Context, productModel model.ProductUpdateModel, id string) model.ProductModel {
	// Validasi input
	common.Validate(productModel)

	// Ambil produk lama dari database
	existingProduct, err := service.ProductRepository.FindById(ctx, id)
	if err != nil {
		panic(exception.NotFoundError{Message: err.Error()})
	}

	// Update field dasar
	if productModel.Name != nil {
		existingProduct.Name = *productModel.Name
	}
	if productModel.Price != nil {
		existingProduct.Price = *productModel.Price
	}
	if productModel.Quantity != nil {
		existingProduct.Quantity = *productModel.Quantity
	}

	// Jika ada image baru
	if productModel.Image != nil {
		fmt.Println("New image uploaded, deleting old and saving new image...")

		// Hapus gambar lama
		service.deleteProductImage(existingProduct.ImagePath)

		// Simpan gambar baru
		imageID := uuid.New().String()
		imageName := fmt.Sprintf("%s.png", imageID)
		imagePath := fmt.Sprintf("./media/products/%s", imageName)

		src, err := productModel.Image.Open()
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

		existingProduct.ImagePath = imagePath
		fmt.Println("New image saved:", imagePath)
	} else {
		fmt.Println("No new image provided")
	}

	// Update ke database
	service.ProductRepository.Update(ctx, existingProduct)

	// Hapus cache Redis agar tidak ambil data lama
	service.Cache.Del(ctx, "product:"+id)

	// Return response
	return model.ProductModel{
		Id:        existingProduct.Id.String(),
		Name:      existingProduct.Name,
		Price:     existingProduct.Price,
		Quantity:  existingProduct.Quantity,
		ImagePath: existingProduct.ImagePath,
	}
}

func (service *productServiceImpl) Delete(ctx context.Context, id string) error {
	product, err := service.ProductRepository.FindById(ctx, id)
	if err != nil {
		panic(exception.NotFoundError{
			Message: err.Error(),
		})
	}

	service.deleteProductImage(product.ImagePath)

	service.ProductRepository.Delete(ctx, product)
	return nil
}

func (service *productServiceImpl) deleteProductImage(imagePath string) {
	if imagePath != "" {
		if err := os.Remove(imagePath); err != nil && !os.IsNotExist(err) {
			panic(exception.InternalServerError{Message: "Failed to delete image"})
		}
	}
}

func (service *productServiceImpl) FindById(ctx context.Context, id string) model.ProductModel {
	productCache := configuration.SetCache[entity.Product](service.Cache, ctx, "product", id, service.ProductRepository.FindById)
	return model.ProductModel{
		Id:        productCache.Id.String(),
		Name:      productCache.Name,
		Price:     productCache.Price,
		Quantity:  productCache.Quantity,
		ImagePath: productCache.ImagePath,
	}
}

func (service *productServiceImpl) FindAll(ctx context.Context) (responses []model.ProductModel) {
	products := service.ProductRepository.FindAl(ctx)
	for _, product := range products {
		responses = append(responses, model.ProductModel{
			Id:        product.Id.String(),
			Name:      product.Name,
			Price:     product.Price,
			Quantity:  product.Quantity,
			ImagePath: product.ImagePath,
		})
	}
	if len(products) == 0 {
		return []model.ProductModel{}
	}
	return responses
}

func (service *productServiceImpl) FindByUsername(ctx context.Context, username string) (responses []model.ProductModel) {
	products := service.ProductRepository.FindByUsername(ctx, username)

	for _, product := range products {
		responses = append(responses, model.ProductModel{
			Id:        product.Id.String(),
			Name:      product.Name,
			Price:     product.Price,
			Quantity:  product.Quantity,
			ImagePath: product.ImagePath,
		})
	}
	if len(products) == 0 {
		return []model.ProductModel{}
	}
	return responses
}
