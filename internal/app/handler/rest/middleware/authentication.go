package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/jwt"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/response"
	"strings"
	"time"
)

type IAuthenticationMiddleware interface {
	Authenticate() fiber.Handler
}

type AuthenticationMiddleware struct {
	jwtAuth jwt.IJWT
}

func NewAuthenticationMiddleware(jwtAuth jwt.IJWT) IAuthenticationMiddleware {
	return AuthenticationMiddleware{jwtAuth: jwtAuth}
}

func (m AuthenticationMiddleware) Authenticate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		bearer := c.Get("Authorization")
		if bearer == "" {
			return response.New(401, "Empty token", nil).Send(c)
		}

		tokenSlice := strings.Split(bearer, " ")
		if len(tokenSlice) != 2 {
			return response.New(401, "Invalid token", nil).Send(c)
		}
		token := tokenSlice[1]
		var claims jwt.Claims
		err := m.jwtAuth.Decode(token, &claims)
		if err != nil {
			return response.New(401, "Fail to validate token", err).Send(c)
		}

		if claims.ExpiresAt.Time.Before(time.Now()) {
			return response.New(401, "Token expired", nil).Send(c)
		}

		c.Locals("claims", claims)
		return c.Next()
	}
}
