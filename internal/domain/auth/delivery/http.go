package delivery

import (
	"github.com/mdanialr/pwman_backend/internal/domain/auth"
	authUC "github.com/mdanialr/pwman_backend/internal/domain/auth/usecase"
	resp "github.com/mdanialr/pwman_backend/pkg/response"

	"github.com/gofiber/fiber/v2"
)

// NewDelivery setup endpoints in domain auth as delivery layer.
func NewDelivery(app fiber.Router, uc authUC.UseCase) {
	d := &delivery{uc: uc}

	api := app.Group("/auth")
	api.Post("/otp", d.LoginOTP)
}

type delivery struct {
	uc authUC.UseCase
}

func (d *delivery) LoginOTP(c *fiber.Ctx) error {
	var req auth.Request
	c.BodyParser(&req)

	if err := req.Validate(); err != nil {
		return resp.Error(c, resp.WithErrValidation(err))
	}

	usr, err := d.uc.ValidateOTP(c.Context(), req)
	if err != nil {
		return resp.Error(c, resp.WithErr(err))
	}

	return resp.Success(c, resp.WithData(usr))
}
