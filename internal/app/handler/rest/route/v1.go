package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/nathakusuma/sea-salon-be/internal/app/handler/rest"
)

type Config struct {
	App           *fiber.App
	ReviewHandler rest.IReviewHandler
}

func (c *Config) Setup() {
	c.App.Use(cors.New())
	v1 := c.App.Group("/v1")

	c.reviewRoutes(v1)
}

func (c *Config) reviewRoutes(r fiber.Router) {
	reviews := r.Group("/reviews")
	reviews.Post("", c.ReviewHandler.Create())
	reviews.Get("", c.ReviewHandler.FindByLazyLoad())
}
