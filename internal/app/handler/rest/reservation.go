package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nathakusuma/sea-salon-be/internal/app/service"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/jwt"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/model"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/response"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/validator"
)

type IReservationHandler interface {
	Create() fiber.Handler
	FindAvailableSchedules() fiber.Handler
	FindByUser() fiber.Handler
}

type reservationHandler struct {
	s service.IReservationService
	v validator.Validator
}

func NewReservationHandler(s service.IReservationService, v validator.Validator) IReservationHandler {
	return &reservationHandler{s: s, v: v}
}

func (h *reservationHandler) Create() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req model.CreateReservationRequest
		if err := c.BodyParser(&req); err != nil {
			return response.New(400, "Fail to parse request body", err.Error()).Send(c)
		}
		if err := h.v.Validate(req); err != nil {
			return response.New(400, "Fail to validate request body", err).Send(c)
		}

		userClaims := c.Locals("claims").(jwt.Claims)

		res := h.s.Create(req, userClaims)
		return res.Send(c)
	}
}

func (h *reservationHandler) FindAvailableSchedules() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req model.FindAvailableReservationSchedulesRequest
		if err := c.QueryParser(&req); err != nil {
			return response.New(400, "Fail to parse request query", err).Send(c)
		}
		if err := h.v.Validate(req); err != nil {
			return response.New(400, "Fail to validate request query", err).Send(c)
		}

		res := h.s.FindAvailableSchedules(req)
		return res.Send(c)
	}
}

func (h *reservationHandler) FindByUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userClaims := c.Locals("claims").(jwt.Claims)

		res := h.s.FindByUser(userClaims)
		return res.Send(c)
	}
}
