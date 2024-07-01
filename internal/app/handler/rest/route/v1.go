package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/nathakusuma/sea-salon-be/internal/app/handler/rest"
	"github.com/nathakusuma/sea-salon-be/internal/app/handler/rest/middleware"
)

type Config struct {
	App                *fiber.App
	AuthenticationMid  middleware.IAuthenticationMiddleware
	ReviewHandler      rest.IReviewHandler
	ReservationHandler rest.IReservationHandler
	AuthHandler        rest.IAuthHandler
	ServiceHandler     rest.IServiceHandler
}

func (c *Config) Setup() {
	c.App.Use(cors.New())
	v1 := c.App.Group("/v1")

	c.reviewRoutes(v1)
	c.reservationRoutes(v1)
	c.authRoutes(v1)
	c.serviceRoutes(v1)
}

func (c *Config) reviewRoutes(r fiber.Router) {
	reviews := r.Group("/reviews")
	reviews.Post("", c.AuthenticationMid.Authenticate(), c.ReviewHandler.Create())
	reviews.Get("", c.ReviewHandler.FindByLazyLoad())
}

func (c *Config) reservationRoutes(r fiber.Router) {
	reservations := r.Group("/reservations")
	reservations.Use(c.AuthenticationMid.Authenticate())
	reservations.Post("", c.ReservationHandler.Create())
	reservations.Get("/available", c.ReservationHandler.FindAvailableSchedules())
	reservations.Get("/my", c.ReservationHandler.FindByUser())
	reservations.Get("/admin", middleware.RequireAdmin(), c.ReservationHandler.FindByDate())
}

func (c *Config) authRoutes(r fiber.Router) {
	auth := r.Group("/auth")
	auth.Post("/register", c.AuthHandler.Register())
	auth.Post("/login", c.AuthHandler.Login())
}

func (c *Config) serviceRoutes(r fiber.Router) {
	services := r.Group("/services")
	services.Post("", c.AuthenticationMid.Authenticate(), middleware.RequireAdmin(), c.ServiceHandler.Create())
	services.Get("", c.ServiceHandler.FindAll())
}
