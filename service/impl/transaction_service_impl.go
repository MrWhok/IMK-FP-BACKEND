package impl

import (
	"context"
	"fmt"
	"strings"

	"github.com/MrWhok/IMK-FP-BACKEND/common"
	"github.com/MrWhok/IMK-FP-BACKEND/entity"
	"github.com/MrWhok/IMK-FP-BACKEND/exception"
	"github.com/MrWhok/IMK-FP-BACKEND/model"
	"github.com/MrWhok/IMK-FP-BACKEND/repository"
	"github.com/MrWhok/IMK-FP-BACKEND/service"
	"github.com/MrWhok/IMK-FP-BACKEND/utils"
	"github.com/google/uuid"
)

func NewTransactionServiceImpl(
	transactionRepository *repository.TransactionRepository,
	cartRepo repository.CartRepository,
	productRepo repository.ProductRepository,
	userRepo repository.UserRepository,
) service.TransactionService {
	return &transactionServiceImpl{
		TransactionRepository: *transactionRepository,
		cartRepo:              cartRepo,
		productRepo:           productRepo,
		userRepo:              userRepo,
	}
}

type transactionServiceImpl struct {
	repository.TransactionRepository
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
	userRepo    repository.UserRepository
}

func (transactionService *transactionServiceImpl) Create(ctx context.Context, transactionModel model.TransactionCreateUpdateModel) model.TransactionCreateUpdateModel {
	common.Validate(transactionModel)
	uuidGenerate := uuid.New()
	var transactionDetails []entity.TransactionDetail
	var totalPrice int64 = 0

	for _, detail := range transactionModel.TransactionDetails {
		totalPrice = totalPrice + detail.SubTotalPrice
		transactionDetails = append(transactionDetails, entity.TransactionDetail{
			TransactionId: uuidGenerate,
			ProductId:     detail.ProductId,
			Id:            uuid.New(),
			SubTotalPrice: detail.SubTotalPrice,
			Price:         detail.Price,
			Quantity:      detail.Quantity,
		})
	}

	transaction := entity.Transaction{
		Id:                 uuidGenerate,
		TotalPrice:         totalPrice,
		UserID:             transactionModel.UserID, // ✅ Make sure it's passed
		Status:             "proses",
		TransactionDetails: transactionDetails,
	}

	transactionService.TransactionRepository.Insert(ctx, transaction)
	return transactionModel
}

func (transactionService *transactionServiceImpl) Delete(ctx context.Context, id string) {
	transaction, err := transactionService.TransactionRepository.FindById(ctx, id)
	if err != nil {
		panic(exception.NotFoundError{
			Message: err.Error(),
		})
	}
	transactionService.TransactionRepository.Delete(ctx, transaction)
}

func (transactionService *transactionServiceImpl) FindById(ctx context.Context, id string) model.TransactionModel {
	transaction, err := transactionService.TransactionRepository.FindById(ctx, id)
	if err != nil {
		panic(exception.NotFoundError{
			Message: err.Error(),
		})
	}
	var transactionDetails []model.TransactionDetailModel
	for _, detail := range transaction.TransactionDetails {
		transactionDetails = append(transactionDetails, model.TransactionDetailModel{
			Id:            detail.Id.String(),
			SubTotalPrice: detail.SubTotalPrice,
			Price:         detail.Price,
			Quantity:      detail.Quantity,
			Product: model.ProductModel{
				Id:       detail.Product.Id.String(),
				Name:     detail.Product.Name,
				Price:    detail.Product.Price,
				Quantity: detail.Product.Quantity,
			},
		})
	}

	return model.TransactionModel{
		Id:                 transaction.Id.String(),
		TotalPrice:         transaction.TotalPrice,
		TransactionDetails: transactionDetails,
	}
}

func (transactionService *transactionServiceImpl) FindAll(ctx context.Context) (responses []model.TransactionModel) {
	transactions := transactionService.TransactionRepository.FindAll(ctx)
	for _, transaction := range transactions {
		var transactionDetails []model.TransactionDetailModel
		for _, detail := range transaction.TransactionDetails {
			transactionDetails = append(transactionDetails, model.TransactionDetailModel{
				Id:            detail.Id.String(),
				SubTotalPrice: detail.SubTotalPrice,
				Price:         detail.Price,
				Quantity:      detail.Quantity,
				Product: model.ProductModel{
					Id:       detail.Product.Id.String(),
					Name:     detail.Product.Name,
					Price:    detail.Product.Price,
					Quantity: detail.Product.Quantity,
				},
			})
		}

		responses = append(responses, model.TransactionModel{
			Id:                 transaction.Id.String(),
			TotalPrice:         transaction.TotalPrice,
			TransactionDetails: transactionDetails,
		})
	}

	return responses
}

func (s *transactionServiceImpl) Checkout(ctx context.Context, username string) []model.TransactionModel {
	// Get cart
	cart, err := s.cartRepo.FindByUsername(ctx, username)
	exception.PanicLogging(err)

	if len(cart.Items) == 0 {
		panic(exception.NotFoundError{Message: "Cart is empty"})
	}

	// Group items by owner
	groupedItems := make(map[string][]entity.CartItem)
	for _, item := range cart.Items {
		product, err := s.productRepo.FindById(ctx, item.ProductID)
		exception.PanicLogging(err)

		groupedItems[product.Owner.Username] = append(groupedItems[product.Owner.Username], item)
	}

	var transactionResponses []model.TransactionModel
	var emailContent strings.Builder

	emailContent.WriteString("Thank you for your purchase!Please contact the seller for more information\n\nHere are your items:\n\n")

	for _, items := range groupedItems {
		transactionId := uuid.New()
		var total int64
		var details []entity.TransactionDetail
		var detailModels []model.TransactionDetailModel

		for _, item := range items {
			subTotal := int64(item.Quantity) * item.Product.Price
			productUUID, err := uuid.Parse(item.ProductID)
			if err != nil {
				panic("invalid ProductID: " + item.ProductID)
			}

			productEntity, err := s.productRepo.FindById(ctx, item.ProductID)
			exception.PanicLogging(err)

			if productEntity.Quantity < item.Quantity {
				panic(exception.BadRequestError{Message: "Insufficient stock for product: " + productEntity.Name})
			}

			productEntity.Quantity -= item.Quantity
			s.productRepo.Update(ctx, productEntity)

			detail := entity.TransactionDetail{
				Id:            uuid.New(),
				TransactionId: transactionId,
				ProductId:     productUUID,
				Price:         item.Product.Price,
				Quantity:      item.Quantity,
				SubTotalPrice: subTotal,
			}
			details = append(details, detail)
			total += subTotal

			// Prepare email line
			emailContent.WriteString(fmt.Sprintf("- %s (Seller: %s, WhatsApp: https://wa.me/%s)\n",
				productEntity.Name, productEntity.Owner.Username, productEntity.Owner.Phone))

			// Prepare response model
			detailModels = append(detailModels, model.TransactionDetailModel{
				Id:            detail.Id.String(),
				SubTotalPrice: detail.SubTotalPrice,
				Price:         detail.Price,
				Quantity:      detail.Quantity,
				Product: model.ProductModel{
					Id:         productEntity.Id.String(),
					Name:       productEntity.Name,
					Price:      productEntity.Price,
					Quantity:   productEntity.Quantity,
					ImagePath:  productEntity.ImagePath,
					Owner:      productEntity.Owner.Username,
					OwnerPhone: productEntity.Owner.Phone,
				},
			})
		}

		transaction := entity.Transaction{
			Id:                 transactionId,
			TotalPrice:         total,
			UserID:             username,
			Status:             "proses",
			TransactionDetails: details,
		}
		s.TransactionRepository.Insert(ctx, transaction)

		transactionResponses = append(transactionResponses, model.TransactionModel{
			Id:                 transactionId.String(),
			TotalPrice:         total,
			Status:             "proses",
			TransactionDetails: detailModels,
		})
	}

	// Add points
	user, err := s.userRepo.FindByUsername(ctx, username)
	exception.PanicLogging(err)
	user.Points += 10
	err = s.userRepo.Update(ctx, user)
	exception.PanicLogging(err)

	// Clear cart
	for _, item := range cart.Items {
		s.cartRepo.DeleteItem(ctx, username, item.ProductID)
	}

	// Send email
	err = utils.SendEmail(user.Email, "Transaction Confirmation", emailContent.String())
	if err != nil {
		fmt.Println("Failed to send email:", err)
	}

	return transactionResponses
}

func (transactionService *transactionServiceImpl) FindByUsername(ctx context.Context, username string) []model.TransactionModel {
	transactions := transactionService.TransactionRepository.FindByUsername(ctx, username)
	var responses []model.TransactionModel

	for _, transaction := range transactions {
		var transactionDetails []model.TransactionDetailModel
		for _, detail := range transaction.TransactionDetails {
			transactionDetails = append(transactionDetails, model.TransactionDetailModel{
				Id:            detail.Id.String(),
				SubTotalPrice: detail.SubTotalPrice,
				Price:         detail.Price,
				Quantity:      detail.Quantity,
				Product: model.ProductModel{
					Id:       detail.Product.Id.String(),
					Name:     detail.Product.Name,
					Price:    detail.Product.Price,
					Quantity: detail.Product.Quantity,
				},
			})
		}

		responses = append(responses, model.TransactionModel{
			Id:                 transaction.Id.String(),
			TotalPrice:         transaction.TotalPrice,
			Status:             transaction.Status,
			TransactionDetails: transactionDetails,
		})
	}

	return responses
}

func (transactionService *transactionServiceImpl) FindByBuyerUsername(ctx context.Context, username string) []model.TransactionModel {
	transactions := transactionService.TransactionRepository.FindByBuyerUsername(ctx, username)
	var responses []model.TransactionModel

	for _, transaction := range transactions {
		var transactionDetails []model.TransactionDetailModel
		for _, detail := range transaction.TransactionDetails {
			transactionDetails = append(transactionDetails, model.TransactionDetailModel{
				Id:            detail.Id.String(),
				SubTotalPrice: detail.SubTotalPrice,
				Price:         detail.Price,
				Quantity:      detail.Quantity,
				Product: model.ProductModel{
					Id:       detail.Product.Id.String(),
					Name:     detail.Product.Name,
					Price:    detail.Product.Price,
					Quantity: detail.Product.Quantity,
				},
			})
		}

		responses = append(responses, model.TransactionModel{
			Id:                 transaction.Id.String(),
			TotalPrice:         transaction.TotalPrice,
			Status:             transaction.Status,
			TransactionDetails: transactionDetails,
		})
	}

	return responses
}

func (transactionService *transactionServiceImpl) UpdateStatus(ctx context.Context, id string, status string) error {
	err := transactionService.TransactionRepository.UpdateStatus(ctx, id, status)
	if err != nil {
		return exception.NotFoundError{Message: "Transaction not found with ID: " + id}
	}
	return nil
}
