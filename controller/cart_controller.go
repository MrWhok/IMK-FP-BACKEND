package controller

import (
	"fmt"

	"github.com/MrWhok/IMK-FP-BACKEND/configuration"
	"github.com/MrWhok/IMK-FP-BACKEND/exception"
	"github.com/MrWhok/IMK-FP-BACKEND/middleware"
	"github.com/MrWhok/IMK-FP-BACKEND/model"
	"github.com/MrWhok/IMK-FP-BACKEND/service"
	"github.com/gofiber/fiber/v2"
)

type CartController struct {
	service.CartService
	configuration.Config
}

func NewCartController(cartService *service.CartService, config configuration.Config) *CartController {
	return &CartController{CartService: *cartService, Config: config}
}

func (cartController CartController) Route(app *fiber.App) {
	cartGroup := app.Group("/v1/api/cart")

	cartGroup.Post("/add", middleware.AuthenticateJWT("user", cartController.Config), cartController.AddToCart)
	cartGroup.Get("/", middleware.AuthenticateJWT("user", cartController.Config), cartController.GetMyCart)
	// cartGroup.Put("/:product_id", middleware.AuthenticateJWT("user", controller.Config), cartController.UpdateCartItem)
	// cartGroup.Delete("/:product_id", middleware.AuthenticateJWT("user", controller.Config), cartController.DeleteCartItem)

}

func (c *CartController) AddToCart(ctx *fiber.Ctx) error {
	fmt.Println("[DEBUG] AddToCart handler called")

	username := ctx.Locals("username").(string)

	var req model.AddToCartRequest
	err := ctx.BodyParser(&req)
	exception.PanicLogging(err)

	response := c.CartService.AddToCart(ctx.Context(), username, req)
	return ctx.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "Product added to cart",
		Data:    response,
	})
}

func (c *CartController) GetMyCart(ctx *fiber.Ctx) error {
	username := ctx.Locals("username").(string)
	cart := c.CartService.GetMyCart(ctx.Context(), username)

	return ctx.JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
		Data:    cart,
	})
}
