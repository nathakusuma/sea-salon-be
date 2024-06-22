package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nathakusuma/sea-salon-be/internal/app/service"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/model"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/response"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/validator"
)

type IReviewHandler interface {
	Create() fiber.Handler
	FindByLazyLoad() fiber.Handler
}

type reviewHandler struct {
	s service.IReviewService
	v validator.Validator
}

func NewReviewHandler(s service.IReviewService, v validator.Validator) IReviewHandler {
	return &reviewHandler{s: s, v: v}
}

func (h *reviewHandler) Create() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req model.CreateReviewRequest
		if err := c.BodyParser(&req); err != nil {
			return response.New(400, "Fail to parse request body", err.Error()).Send(c)
		}
		if err := h.v.Validate(req); err != nil {
			return response.New(400, "Fail to validate request body", err).Send(c)
		}

		res := h.s.Create(req)
		return res.Send(c)
	}
}

func (h *reviewHandler) FindByLazyLoad() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req model.FindReviewsLazyLoadRequest
		if err := c.QueryParser(&req); err != nil {
			return response.New(400, "Fail to parse request query", err).Send(c)
		}
		if err := h.v.Validate(req); err != nil {
			return response.New(400, "Fail to validate request query", err).Send(c)
		}

		res := h.s.FindByLazyLoad(req)
		return res.Send(c)
	}
}
