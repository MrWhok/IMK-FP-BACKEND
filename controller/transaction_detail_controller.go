package controller

import (
	"github.com/MrWhok/IMK-FP-BACKEND/configuration"
	"github.com/MrWhok/IMK-FP-BACKEND/middleware"
	"github.com/MrWhok/IMK-FP-BACKEND/model"
	"github.com/MrWhok/IMK-FP-BACKEND/service"
	"github.com/gofiber/fiber/v2"
)

type TransactionDetailController struct {
	service.TransactionDetailService
	configuration.Config
}

func NewTransactionDetailController(transactionDetailService *service.TransactionDetailService, config configuration.Config) *TransactionDetailController {
	return &TransactionDetailController{TransactionDetailService: *transactionDetailService, Config: config}
}

func (controller TransactionDetailController) Route(app *fiber.App) {
	app.Get("/v1/api/transaction-detail/:id", middleware.AuthenticateJWT("user", controller.Config), controller.FindById)
}

func (controller TransactionDetailController) FindById(c *fiber.Ctx) error {
	id := c.Params("id")

	result := controller.TransactionDetailService.FindById(c.Context(), id)
	return c.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
		Data:    result,
	})
}
