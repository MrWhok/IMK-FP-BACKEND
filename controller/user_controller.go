package controller

import (
	"github.com/MrWhok/IMK-FP-BACKEND/common"
	"github.com/MrWhok/IMK-FP-BACKEND/configuration"
	"github.com/MrWhok/IMK-FP-BACKEND/exception"
	"github.com/MrWhok/IMK-FP-BACKEND/middleware"
	"github.com/MrWhok/IMK-FP-BACKEND/model"
	"github.com/MrWhok/IMK-FP-BACKEND/service"
	"github.com/gofiber/fiber/v2"
)

func NewUserController(userService *service.UserService, config configuration.Config) *UserController {
	return &UserController{UserService: *userService, Config: config}
}

type UserController struct {
	service.UserService
	configuration.Config
}

func (controller UserController) Route(app *fiber.App) {
	app.Post("/v1/api/authentication", controller.Authentication)
	app.Post("/v1/api/register", controller.Register)
	app.Get("/v1/api/me", middleware.AuthenticateJWT("user", controller.Config), controller.Me)
	app.Get("/v1/api/leaderboard", middleware.AuthenticateJWT("user", controller.Config), controller.Leaderboard)
	app.Put("/v1/api/update-profile", middleware.AuthenticateJWT("user", controller.Config), controller.UpdateProfile)
}

// Authentication func Authenticate user.
// @Description authenticate user.
// @Summary authenticate user
// @Tags Authenticate user
// @Accept json
// @Produce json
// @Param request body model.UserModel true "Request Body"
// @Success 200 {object} model.GeneralResponse
// @Router /v1/api/authentication [post]
func (controller UserController) Authentication(c *fiber.Ctx) error {
	var request model.UserModel
	err := c.BodyParser(&request)
	exception.PanicLogging(err)

	result := controller.UserService.Authentication(c.Context(), request)
	var userRoles []map[string]interface{}
	for _, userRole := range result.UserRoles {
		userRoles = append(userRoles, map[string]interface{}{
			"role": userRole.Role,
		})
	}
	tokenJwtResult, expirationTime := common.GenerateToken(result.Username, userRoles, controller.Config)

	cookie := new(fiber.Cookie)
	cookie.Name = "token"
	cookie.Value = tokenJwtResult
	cookie.Expires = expirationTime
	cookie.HTTPOnly = true
	cookie.Secure = true
	cookie.SameSite = fiber.CookieSameSiteNoneMode
	cookie.Path = "/"

	c.Cookie(cookie)

	resultWithToken := map[string]interface{}{
		"token":    tokenJwtResult,
		"username": result.Username,
		"role":     userRoles,
	}
	return c.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
		Data:    resultWithToken,
	})
}

func (controller UserController) Register(c *fiber.Ctx) error {
	var request model.UserCreateModel
	err := c.BodyParser(&request)
	exception.PanicLogging(err)

	err = controller.UserService.Register(c.Context(), request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.GeneralResponse{
			Code:    400,
			Message: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    201,
		Message: "User registered successfully",
	})
}

func (controller UserController) Me(c *fiber.Ctx) error {
	usernameInterface := c.Locals("username")
	username, ok := usernameInterface.(string)
	if !ok || username == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(model.GeneralResponse{
			Code:    401,
			Message: "Unauthorized",
			Data:    "Invalid or missing user",
		})
	}

	response, err := controller.UserService.FindMe(c.Context(), username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.GeneralResponse{
			Code:    500,
			Message: "General Error",
			Data:    err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
		Data:    response,
	})
}

func (controller UserController) Leaderboard(c *fiber.Ctx) error {
	leaderboard, err := controller.UserService.GetLeaderboard(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.GeneralResponse{
			Code:    500,
			Message: "Failed to fetch leaderboard",
			Data:    err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
		Data:    leaderboard,
	})
}

func (controller UserController) UpdateProfile(c *fiber.Ctx) error {
	var request model.UserUpdateModel
	err := c.BodyParser(&request)
	exception.PanicLogging(err)

	usernameInterface := c.Locals("username")
	username, ok := usernameInterface.(string)
	if !ok || username == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(model.GeneralResponse{
			Code:    401,
			Message: "Unauthorized",
			Data:    "Invalid or missing user",
		})
	}

	err = controller.UserService.UpdateProfile(c.Context(), username, request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.GeneralResponse{
			Code:    400,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "Profile updated successfully",
	})
}
