package response

import (
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	HttpCode int    `json:"-"`
	Message  string `json:"message"`
	Payload  any    `json:"payload"`
}

func New(httpCode int, message string, payload any) Response {
	if payload == nil {
		payload = map[string]any{}
	}
	return Response{
		HttpCode: httpCode,
		Message:  message,
		Payload:  payload,
	}
}

func (r Response) Send(c *fiber.Ctx) error {
	return c.Status(r.HttpCode).JSON(r)
}
