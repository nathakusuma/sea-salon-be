package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nathakusuma/sea-salon-be/internal/app/handler/rest"
	"github.com/nathakusuma/sea-salon-be/internal/app/handler/rest/route"
	"github.com/nathakusuma/sea-salon-be/internal/app/repository"
	"github.com/nathakusuma/sea-salon-be/internal/app/service"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/validator"
	"gorm.io/gorm"
)

type StartAppConfig struct {
	DB  *gorm.DB
	App *fiber.App
}

func StartApp(config *StartAppConfig) {
	val := validator.NewValidator()

	// Repositories
	reviewRepo := repository.NewReviewRepository(config.DB)
	reservationRepo := repository.NewReservationRepository(config.DB)

	// Services
	reviewService := service.NewReviewService(reviewRepo)
	reservationService := service.NewReservationService(reservationRepo)

	// Handlers
	reviewHandler := rest.NewReviewHandler(reviewService, val)
	reservationHandler := rest.NewReservationHandler(reservationService, val)

	routeConfig := route.Config{
		App:                config.App,
		ReviewHandler:      reviewHandler,
		ReservationHandler: reservationHandler,
	}
	routeConfig.Setup()
}
