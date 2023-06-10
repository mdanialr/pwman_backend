package app

import (
	auth "github.com/mdanialr/pwman_backend/internal/domain/auth/delivery"
	authRepo "github.com/mdanialr/pwman_backend/internal/domain/auth/repository"
	authUC "github.com/mdanialr/pwman_backend/internal/domain/auth/usecase"
	pw "github.com/mdanialr/pwman_backend/internal/domain/password/delivery"
	pwRepo "github.com/mdanialr/pwman_backend/internal/domain/password/repository"
	pwUC "github.com/mdanialr/pwman_backend/internal/domain/password/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// HttpHandler handler that use HTTP in the delivery layer.
type HttpHandler struct {
	R      fiber.Router
	Log    *zap.Logger
	DB     *gorm.DB
	Config *viper.Viper
}

// SetupRouter init all HTTP endpoints and their dependencies.
func (h *HttpHandler) SetupRouter() {
	h.R.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("OK")
	})
	h.R.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("Password Manager API")
	})

	// currently use v1
	v1 := h.R.Group("/v1")

	// init repositories
	authRepository := authRepo.NewRepository(h.DB)
	pwRepository := pwRepo.NewRepository(h.DB)

	// init use cases
	authUseCase := authUC.NewUseCase(h.Config, h.Log, authRepository)
	pwUseCase := pwUC.NewUseCase(h.Config, h.Log, pwRepository)

	// init handlers
	auth.NewDelivery(v1, authUseCase)       // - /auth/*
	pw.NewDelivery(v1, h.Config, pwUseCase) // - /category/*
}
