package controller

import (
	"github.com/MrWhok/IMK-FP-BACKEND/service"
	"github.com/gofiber/fiber/v2"
)

type NewsController struct {
	NewsService service.NewsService
}

func NewNewsController(newsService *service.NewsService) *NewsController {
	return &NewsController{NewsService: *newsService}
}

func (controller NewsController) Route(app *fiber.App) {
	newsGroup := app.Group("/v1/api/news")

	newsGroup.Get("", controller.FindAll)
	newsGroup.Get("/status", controller.Status)
}

func (controller *NewsController) FindAll(c *fiber.Ctx) error {
	return c.JSON(controller.NewsService.GetNews())
}

func (controller *NewsController) Status(c *fiber.Ctx) error {
	data := controller.NewsService.GetNews()
	return c.JSON(fiber.Map{
		"lastUpdated": data.LastUpdated,
		"count":       len(data.Data),
	})
}
