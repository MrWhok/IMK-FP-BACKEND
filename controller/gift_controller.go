package controller

import (
	"fmt"
	"strconv"

	"github.com/MrWhok/IMK-FP-BACKEND/configuration"
	"github.com/MrWhok/IMK-FP-BACKEND/exception"
	"github.com/MrWhok/IMK-FP-BACKEND/middleware"
	"github.com/MrWhok/IMK-FP-BACKEND/model"
	"github.com/MrWhok/IMK-FP-BACKEND/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type GiftController struct {
	service.GiftService
	configuration.Config
}

func NewGiftController(giftService *service.GiftService, config configuration.Config) *GiftController {
	return &GiftController{GiftService: *giftService, Config: config}
}

func (controller GiftController) Route(app *fiber.App) {
	giftGroup := app.Group("/v1/api/gift")

	giftGroup.Post("", middleware.AuthenticateJWT("admin", controller.Config), controller.Create)
	giftGroup.Put("/:id", middleware.AuthenticateJWT("admin", controller.Config), controller.Update)
	giftGroup.Delete("/:id", middleware.AuthenticateJWT("admin", controller.Config), controller.Delete)
	giftGroup.Get("/:id", middleware.AuthenticateJWT("user", controller.Config), controller.FindById)
	giftGroup.Get("", middleware.AuthenticateJWT("user", controller.Config), controller.FindAll)
	giftGroup.Post("/:giftId/exchange", middleware.AuthenticateJWT("user", controller.Config), controller.ExchangeGift)
}

func (controller GiftController) Create(c *fiber.Ctx) error {
	// Parse multipart form fields manually
	name := c.FormValue("name")
	priceStr := c.FormValue("point_price")
	quantityStr := c.FormValue("quantity")

	// fmt.Println("DEBUG FORM:", name, priceStr, quantityStr, category)

	if name == "" || priceStr == "" || quantityStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.GeneralResponse{
			Code:    400,
			Message: "Missing required fields",
			Data: []map[string]string{
				{"field": "Name", "message": "this field is required"},
				{"field": "Price", "message": "this field is required"},
				{"field": "Quantity", "message": "this field is required"},
			},
		})
	}

	// Convert price and quantity
	price, err := strconv.ParseInt(priceStr, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid price format")
	}

	quantity, err := strconv.ParseInt(quantityStr, 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid quantity format")
	}

	// Get file
	file, err := c.FormFile("image")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Image is required")
	}

	// Save image to local folder
	imageID := uuid.New().String()
	imageName := fmt.Sprintf("%s.png", imageID)

	imagePath := fmt.Sprintf("./media/gifts/%s", imageName)
	if err := c.SaveFile(file, imagePath); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to save image")
	}

	// Create request object
	request := model.GiftCreateModel{
		Name:       name,
		PointPrice: price,
		Quantity:   int32(quantity),
		Image:      file,
	}

	// Call service
	response := controller.GiftService.Create(c.Context(), request, imagePath)

	return c.Status(fiber.StatusCreated).JSON(model.GeneralResponse{
		Code:    200,
		Message: "Gift created successfully",
		Data:    response,
	})
}

func (controller GiftController) Update(c *fiber.Ctx) error {
	var request model.GiftUpdateModel
	id := c.Params("id")
	err := c.BodyParser(&request)
	exception.PanicLogging(err)

	file, err := c.FormFile("image")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Image is required")
	}
	request.Image = file

	response := controller.GiftService.Update(c.Context(), request, id)
	return c.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
		Data:    response,
	})
}

func (controller GiftController) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	controller.GiftService.Delete(c.Context(), id)
	return c.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
	})
}

func (controller GiftController) FindById(c *fiber.Ctx) error {
	id := c.Params("id")

	result := controller.GiftService.FindById(c.Context(), id)
	return c.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
		Data:    result,
	})
}

func (controller GiftController) FindAll(c *fiber.Ctx) error {
	result := controller.GiftService.FindAll(c.Context())
	fmt.Println("In the findAll controller:")
	return c.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
		Data:    result,
	})
}

func (controller GiftController) ExchangeGift(c *fiber.Ctx) error {
	giftId := c.Params("giftId")
	username := c.Locals("username").(string)

	err := controller.GiftService.ExchangeGift(c.Context(), giftId, username)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.GeneralResponse{
			Code:    400,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "Gift exchanged successfully",
	})
}
