package response

import (
	"github.com/gofiber/fiber/v2"
	"log"
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
	if r.HttpCode >= 500 {
		log.Printf("ERROR %v: %v payload=%v\n", r.HttpCode, r.Message, r.Payload)
		r.Payload = map[string]any{}
	}
	return c.Status(r.HttpCode).JSON(r)
}
