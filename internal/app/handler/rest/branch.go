package rest

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/nathakusuma/sea-salon-be/internal/app/service"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/model"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/response"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/validator"
)

type IBranchHandler interface {
	Create() fiber.Handler
	FindAll() fiber.Handler
	SetServices() fiber.Handler
}

type branchHandler struct {
	s service.IBranchService
	v validator.Validator
}

func NewBranchHandler(s service.IBranchService, v validator.Validator) IBranchHandler {
	return &branchHandler{s: s, v: v}
}

func (h *branchHandler) Create() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req model.CreateBranchRequest

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

func (h *branchHandler) FindAll() fiber.Handler {
	return func(c *fiber.Ctx) error {
		res := h.s.FindAll()
		return res.Send(c)
	}
}

func (h *branchHandler) SetServices() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req model.SetServicesToBranchRequest
		if err := c.BodyParser(&req); err != nil {
			return response.New(400, "Fail to parse request body", err.Error()).Send(c)
		}
		if err := h.v.Validate(req); err != nil {
			return response.New(400, "Fail to validate request body", err).Send(c)
		}

		res := h.s.SetServices(req)
		return res.Send(c)
	}
}
