package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nathakusuma/sea-salon-be/internal/app/service"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/model"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/response"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/validator"
)

type IAuthHandler interface {
	Register() fiber.Handler
	Login() fiber.Handler
}

type authHandler struct {
	s service.IAuthService
	v validator.Validator
}

func NewAuthHandler(s service.IAuthService, v validator.Validator) IAuthHandler {
	return &authHandler{s: s, v: v}
}

func (h *authHandler) Register() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req model.RegisterRequest
		if err := c.BodyParser(&req); err != nil {
			return response.New(400, "Fail to parse request body", err.Error()).Send(c)
		}
		if err := h.v.Validate(req); err != nil {
			return response.New(400, "Fail to validate request body", err).Send(c)
		}

		res := h.s.Register(req)
		return res.Send(c)
	}
}

func (h *authHandler) Login() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req model.LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return response.New(400, "Fail to parse request body", err.Error()).Send(c)
		}
		if err := h.v.Validate(req); err != nil {
			return response.New(400, "Fail to validate request body", err).Send(c)
		}

		res := h.s.Login(req)
		return res.Send(c)
	}
}
