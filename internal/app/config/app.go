package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nathakusuma/sea-salon-be/internal/app/handler/rest"
	"github.com/nathakusuma/sea-salon-be/internal/app/handler/rest/middleware"
	"github.com/nathakusuma/sea-salon-be/internal/app/handler/rest/route"
	"github.com/nathakusuma/sea-salon-be/internal/app/repository"
	"github.com/nathakusuma/sea-salon-be/internal/app/service"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/jwt"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/validator"
	"gorm.io/gorm"
	"os"
)

type StartAppConfig struct {
	DB  *gorm.DB
	App *fiber.App
}

func StartApp(config *StartAppConfig) {
	val := validator.NewValidator()
	jwtAuth := jwt.NewJWT(os.Getenv("JWT_SECRET_KEY"), os.Getenv("JWT_TTL"))

	// Repositories
	reviewRepo := repository.NewReviewRepository(config.DB)
	reservationRepo := repository.NewReservationRepository(config.DB)
	userRepo := repository.NewUserRepository(config.DB)

	// Services
	reviewService := service.NewReviewService(reviewRepo)
	reservationService := service.NewReservationService(reservationRepo)
	authService := service.NewAuthService(userRepo, jwtAuth)

	// Middlewares
	authenticationMid := middleware.NewAuthenticationMiddleware(jwtAuth)

	// Handlers
	reviewHandler := rest.NewReviewHandler(reviewService, val)
	reservationHandler := rest.NewReservationHandler(reservationService, val)
	authHandler := rest.NewAuthHandler(authService, val)

	routeConfig := route.Config{
		App:                config.App,
		ReviewHandler:      reviewHandler,
		ReservationHandler: reservationHandler,
		AuthHandler:        authHandler,
		AuthenticationMid:  authenticationMid,
	}
	routeConfig.Setup()
}
