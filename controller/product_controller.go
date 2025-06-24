package controller

import (
	"fmt"
	"strconv"

	"github.com/MrWhok/IMK-FP-BACKEND/configuration"
	"github.com/MrWhok/IMK-FP-BACKEND/middleware"
	"github.com/MrWhok/IMK-FP-BACKEND/model"
	"github.com/MrWhok/IMK-FP-BACKEND/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProductController struct {
	service.ProductService
	configuration.Config
}

func NewProductController(productService *service.ProductService, config configuration.Config) *ProductController {
	return &ProductController{ProductService: *productService, Config: config}
}

func (controller ProductController) Route(app *fiber.App) {
	productGroup := app.Group("/v1/api/product")

	productGroup.Get("/myproducts", middleware.AuthenticateJWT("user", controller.Config), controller.MyProducts)
	productGroup.Post("", middleware.AuthenticateJWT("user", controller.Config), controller.Create)
	productGroup.Put("/:id", middleware.AuthenticateJWT("user", controller.Config), controller.Update)
	productGroup.Delete("/:id", middleware.AuthenticateJWT("user", controller.Config), controller.Delete)
	productGroup.Get("/:id", middleware.AuthenticateJWT("user", controller.Config), controller.FindById)
	productGroup.Get("", middleware.AuthenticateJWT("user", controller.Config), controller.FindAll)
}

// Create func create product.
// @Description create product.
// @Summary create product
// @Tags Product
// @Accept json
// @Produce json
// @Param request body model.ProductCreateOrUpdateModel true "Request Body"
// @Success 200 {object} model.GeneralResponse
// @Security JWT
// @Router /v1/api/product [post]
func (controller ProductController) Create(c *fiber.Ctx) error {
	// Parse multipart form fields manually
	name := c.FormValue("name")
	priceStr := c.FormValue("price")
	quantityStr := c.FormValue("quantity")
	category := c.FormValue("category")
	description := c.FormValue("description")

	fmt.Println("DEBUG FORM:", name, priceStr, quantityStr, category)

	if name == "" || priceStr == "" || quantityStr == "" || category == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.GeneralResponse{
			Code:    400,
			Message: "Missing required fields",
			Data: []map[string]string{
				{"field": "Name", "message": "this field is required"},
				{"field": "Price", "message": "this field is required"},
				{"field": "Quantity", "message": "this field is required"},
				{"field": "Category", "message": "this field is required"},
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

	imagePath := fmt.Sprintf("./media/products/%s", imageName)
	if err := c.SaveFile(file, imagePath); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to save image")
	}

	// Create request object
	request := model.ProductCreateModel{
		Name:        name,
		Price:       price,
		Quantity:    int32(quantity),
		Category:    category,
		Description: description,
		Image:       file,
	}

	// Call service
	response := controller.ProductService.Create(c.Context(), request, imagePath)

	return c.Status(fiber.StatusCreated).JSON(model.GeneralResponse{
		Code:    200,
		Message: "Product created successfully",
		Data:    response,
	})
}

// Update func update one exists product.
// @Description update one exists product.
// @Summary update one exists product
// @Tags Product
// @Accept json
// @Produce json
// @Param request body model.ProductCreateOrUpdateModel true "Request Body"
// @Param id path string true "Product Id"
// @Success 200 {object} model.GeneralResponse
// @Security JWT
// @Router /v1/api/product/{id} [put]
func (controller ProductController) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	// err := c.BodyParser(&request)
	// exception.PanicLogging(err)
	name := c.FormValue("name")
	priceStr := c.FormValue("price")
	quantityStr := c.FormValue("quantity")
	category := c.FormValue("category")
	description := c.FormValue("description")

	if name == "" || priceStr == "" || quantityStr == "" || category == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.GeneralResponse{
			Code:    400,
			Message: "Missing required fields",
			Data: []map[string]string{
				{"field": "Name", "message": "this field is required"},
				{"field": "Price", "message": "this field is required"},
				{"field": "Quantity", "message": "this field is required"},
				{"field": "Category", "message": "this field is required"},
			},
		})
	}

	price, err := strconv.ParseInt(priceStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.GeneralResponse{
			Code: 400, Message: "Format harga tidak valid.",
		})
	}

	quantity64, err := strconv.ParseInt(quantityStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.GeneralResponse{
			Code: 400, Message: "Format jumlah tidak valid.",
		})
	}

	request := model.ProductUpdateModel{
		Name:        name,
		Price:       price,
		Quantity:    int32(quantity64),
		Description: description,
		Category:    category,
	}

	file, _ := c.FormFile("image")
	// if err != nil {
	// 	return fiber.NewError(fiber.StatusBadRequest, "Image is required")
	// }
	request.Image = file

	response := controller.ProductService.Update(c.Context(), request, id)
	return c.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
		Data:    response,
	})
}

// Delete func delete one exists product.
// @Description delete one exists product.
// @Summary delete one exists product
// @Tags Product
// @Accept json
// @Produce json
// @Param id path string true "Product Id"
// @Success 200 {object} model.GeneralResponse
// @Security JWT
// @Router /v1/api/product/{id} [delete]
func (controller ProductController) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	controller.ProductService.Delete(c.Context(), id)
	return c.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
	})
}

// FindById func gets one exists product.
// @Description Get one exists product.
// @Summary get one exists product
// @Tags Product
// @Accept json
// @Produce json
// @Param id path string true "Product Id"
// @Success 200 {object} model.GeneralResponse
// @Security JWT
// @Router /v1/api/product/{id} [get]
func (controller ProductController) FindById(c *fiber.Ctx) error {
	id := c.Params("id")

	result := controller.ProductService.FindById(c.Context(), id)
	return c.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
		Data:    result,
	})
}

// FindAll func gets all exists products.
// @Description Get all exists products.
// @Summary get all exists products
// @Tags Product
// @Accept json
// @Produce json
// @Success 200 {object} model.GeneralResponse
// @Security JWT
// @Router /v1/api/product [get]
func (controller ProductController) FindAll(c *fiber.Ctx) error {
	result := controller.ProductService.FindAll(c.Context())
	fmt.Println("In the findAll controller:")
	return c.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
		Data:    result,
	})
}

func (controller ProductController) MyProducts(c *fiber.Ctx) error {
	username := c.Locals("username").(string)

	fmt.Println("Username from JWT:", username)

	result := controller.ProductService.FindByUsername(c.Context(), username)
	return c.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
		Data:    result,
	})
}
