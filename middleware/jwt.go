package middleware

import (
	"github.com/MrWhok/IMK-FP-BACKEND/common"
	"github.com/MrWhok/IMK-FP-BACKEND/configuration"
	"github.com/MrWhok/IMK-FP-BACKEND/model"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
)

func AuthenticateJWT(role string, config configuration.Config) func(*fiber.Ctx) error {
	jwtSecret := config.Get("JWT_SECRET_KEY")
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(jwtSecret),
		SuccessHandler: func(ctx *fiber.Ctx) error {
			user := ctx.Locals("user").(*jwt.Token)
			claims := user.Claims.(jwt.MapClaims)

			if username, ok := claims["username"].(string); ok {
				ctx.Locals("username", username)
			} else {
				return ctx.
					Status(fiber.StatusUnauthorized).
					JSON(model.GeneralResponse{
						Code:    401,
						Message: "Unauthorized",
						Data:    "Username claim missing",
					})
			}

			roles := claims["roles"].([]interface{})

			common.NewLogger().Info("role function ", role, " role user ", roles)
			for _, roleInterface := range roles {
				roleMap := roleInterface.(map[string]interface{})
				if roleMap["role"] == role {
					return ctx.Next()
				}
			}

			return ctx.
				Status(fiber.StatusUnauthorized).
				JSON(model.GeneralResponse{
					Code:    401,
					Message: "Unauthorized",
					Data:    "Invalid Role",
				})
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if err.Error() == "Missing or malformed JWT" {
				return c.
					Status(fiber.StatusBadRequest).
					JSON(model.GeneralResponse{
						Code:    400,
						Message: "Bad Request",
						Data:    "Missing or malformed JWT",
					})
			} else {
				return c.
					Status(fiber.StatusUnauthorized).
					JSON(model.GeneralResponse{
						Code:    401,
						Message: "Unauthorized",
						Data:    "Invalid or expired JWT",
					})
			}
		},
	})
}

func JWTAuth(config configuration.Config) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(config.Get("JWT_SECRET_KEY")),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(model.GeneralResponse{
				Code:    401,
				Message: "Unauthorized",
				Data:    "Invalid or expired JWT",
			})
		},
	})
}

func ExtractUserContext() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)

		username, ok := claims["username"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(model.GeneralResponse{
				Code:    401,
				Message: "Unauthorized",
				Data:    "Username claim missing",
			})
		}
		c.Locals("username", username)
		c.Locals("roles", claims["roles"]) // optional, for next middleware
		return c.Next()
	}
}

func AuthorizeRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rolesRaw := c.Locals("roles")
		roles, ok := rolesRaw.([]interface{})
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(model.GeneralResponse{
				Code:    401,
				Message: "Unauthorized",
				Data:    "Roles claim missing or invalid",
			})
		}

		for _, r := range roles {
			roleMap := r.(map[string]interface{})
			if roleMap["role"] == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusUnauthorized).JSON(model.GeneralResponse{
			Code:    401,
			Message: "Unauthorized",
			Data:    "Insufficient permissions",
		})
	}
}
