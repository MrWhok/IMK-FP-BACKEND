package configuration

import (
	"github.com/MrWhok/IMK-FP-BACKEND/exception"
	"github.com/gofiber/fiber/v2"
)

func NewFiberConfiguration() fiber.Config {
	return fiber.Config{
		ErrorHandler: exception.ErrorHandler,
	}
}
