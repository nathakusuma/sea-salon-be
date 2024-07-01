package rest

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/nathakusuma/sea-salon-be/internal/app/service"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/model"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/response"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/validator"
)

type IServiceHandler interface {
	Create() fiber.Handler
	FindAll() fiber.Handler
}

type serviceHandler struct {
	s service.IServiceService
	v validator.Validator
}

func NewServiceHandler(s service.IServiceService, v validator.Validator) IServiceHandler {
	return &serviceHandler{s: s, v: v}
}

func (h *serviceHandler) Create() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req model.CreateServiceRequest

		imageFileHeader, err := c.FormFile("image")
		if err != nil {
			return response.New(400, "Fail to get image", err.Error()).Send(c)
		}
		req.ImageFile = imageFileHeader

		jsonStr := c.FormValue("json")
		if err := json.Unmarshal([]byte(jsonStr), &req); err != nil {
			return response.New(400, "Fail to unmarshal request body", err).Send(c)
		}

		if err := h.v.Validate(req); err != nil {
			return response.New(400, "Fail to validate request body", err).Send(c)
		}

		res := h.s.Create(&req)
		return res.Send(c)
	}
}

func (h *serviceHandler) FindAll() fiber.Handler {
	return func(c *fiber.Ctx) error {
		res := h.s.FindAll()
		return res.Send(c)
	}
}
