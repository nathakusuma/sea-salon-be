package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/jwt"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/response"
)

func RequireAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims := c.Locals("claims").(jwt.Claims)
		if !claims.IsAdmin {
			return response.New(403, "You're not an admin", nil).Send(c)
		}

		return c.Next()
	}
}
