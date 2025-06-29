package controller

import (
	"github.com/MrWhok/IMK-FP-BACKEND/configuration"
	"github.com/MrWhok/IMK-FP-BACKEND/exception"
	"github.com/MrWhok/IMK-FP-BACKEND/middleware"
	"github.com/MrWhok/IMK-FP-BACKEND/model"
	"github.com/MrWhok/IMK-FP-BACKEND/service"
	"github.com/gofiber/fiber/v2"
)

type TransactionController struct {
	service.TransactionService
	configuration.Config
}

func NewTransactionController(transactionService *service.TransactionService, config configuration.Config) *TransactionController {
	return &TransactionController{TransactionService: *transactionService, Config: config}
}

func (controller TransactionController) Route(app *fiber.App) {
	transactionGroup := app.Group("/v1/api/transaction")

	transactionGroup.Post("", middleware.AuthenticateJWT("user", controller.Config), controller.Create)
	transactionGroup.Get("/my", middleware.AuthenticateJWT("user", controller.Config), controller.FindByUsername)
	transactionGroup.Get("/buyer", middleware.AuthenticateJWT("user", controller.Config), controller.FindByBuyerUsername)
	transactionGroup.Delete("/:id", middleware.AuthenticateJWT("user", controller.Config), controller.Delete)
	transactionGroup.Get("/:id", middleware.AuthenticateJWT("user", controller.Config), controller.FindById)
	transactionGroup.Put("/:id", middleware.AuthenticateJWT("user", controller.Config), controller.UpdateStatus)
	transactionGroup.Get("", middleware.AuthenticateJWT("user", controller.Config), controller.FindAll)
	transactionGroup.Post("/checkout", middleware.AuthenticateJWT("user", controller.Config), controller.Checkout)
}

// Create func create transaction.
// @Description create transaction.
// @Summary create transaction
// @Tags Transaction
// @Accept json
// @Produce json
// @Param request body model.TransactionCreateUpdateModel true "Request Body"
// @Success 200 {object} model.GeneralResponse
// @Security JWT
// @Router /v1/api/transaction [post]
func (controller TransactionController) Create(c *fiber.Ctx) error {
	var request model.TransactionCreateUpdateModel
	err := c.BodyParser(&request)
	exception.PanicLogging(err)

	username := c.Locals("username").(string) // ✅ Get username from JWT
	request.UserID = username                 // ✅ Ensure correct user is stored

	response := controller.TransactionService.Create(c.Context(), request)
	return c.Status(fiber.StatusCreated).JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
		Data:    response,
	})
}

// Delete func delete one exists transaction.
// @Description delete one exists transaction.
// @Summary delete one exists transaction
// @Tags Transaction
// @Accept json
// @Produce json
// @Param id path string true "Transaction Id"
// @Success 200 {object} model.GeneralResponse
// @Security JWT
// @Router /v1/api/transaction/{id} [delete]
func (controller TransactionController) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	controller.TransactionService.Delete(c.Context(), id)
	return c.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
	})
}

// FindById func gets one exists transaction.
// @Description Get one exists transaction.
// @Summary get one exists transaction
// @Tags Transaction
// @Accept json
// @Produce json
// @Param id path string true "Transaction Id"
// @Success 200 {object} model.GeneralResponse
// @Security JWT
// @Router /v1/api/transaction/{id} [get]
func (controller TransactionController) FindById(c *fiber.Ctx) error {
	id := c.Params("id")

	result := controller.TransactionService.FindById(c.Context(), id)
	return c.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
		Data:    result,
	})
}

// FindAll func gets all exists transaction.
// @Description Get all exists transaction.
// @Summary get all exists transaction
// @Tags Transaction
// @Accept json
// @Produce json
// @Success 200 {object} model.GeneralResponse
// @Security JWT
// @Router /v1/api/transaction [get]
func (controller TransactionController) FindAll(c *fiber.Ctx) error {
	result := controller.TransactionService.FindAll(c.Context())
	return c.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
		Data:    result,
	})
}

func (controller TransactionController) Checkout(ctx *fiber.Ctx) error {
	username := ctx.Locals("username").(string)
	result := controller.TransactionService.Checkout(ctx.Context(), username)
	return ctx.JSON(model.GeneralResponse{
		Code:    200,
		Message: "Checkout success",
		Data:    result,
	})
}

func (controller TransactionController) FindByUsername(ctx *fiber.Ctx) error {
	username := ctx.Locals("username").(string)
	result := controller.TransactionService.FindByUsername(ctx.Context(), username)
	return ctx.JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
		Data:    result,
	})
}

func (controller TransactionController) FindByBuyerUsername(ctx *fiber.Ctx) error {
	username := ctx.Locals("username").(string)
	result := controller.TransactionService.FindByBuyerUsername(ctx.Context(), username)
	return ctx.JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
		Data:    result,
	})
}

func (controller TransactionController) UpdateStatus(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	status := ctx.Query("status")

	if status == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.GeneralResponse{
			Code:    400,
			Message: "Status is required",
		})
	}

	err := controller.TransactionService.UpdateStatus(ctx.Context(), id, status)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(model.GeneralResponse{
			Code:    404,
			Message: err.Error(),
		})
	}

	return ctx.JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
	})
}
