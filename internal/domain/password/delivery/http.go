package delivery

import (
	pw "github.com/mdanialr/pwman_backend/internal/domain/password"
	pwUC "github.com/mdanialr/pwman_backend/internal/domain/password/usecase"
	md "github.com/mdanialr/pwman_backend/internal/middleware"
	resp "github.com/mdanialr/pwman_backend/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

// NewDelivery setup endpoints in domain password as delivery layer.
func NewDelivery(app fiber.Router, conf *viper.Viper, uc pwUC.UseCase) {
	d := &delivery{uc: uc}

	api := app.Group("/category", md.JWT(conf))
	api.Get("/", d.Index)
}

type delivery struct {
	uc pwUC.UseCase
}

func (d *delivery) Index(c *fiber.Ctx) error {
	var req pw.Request
	c.QueryParser(&req)
	// set up the query order and sort
	req.SetQuery()

	res, err := d.uc.IndexCategory(c.Context(), req)
	if err != nil {
		return resp.Error(c, resp.WithErr(err))
	}

	return resp.Success(c, resp.WithData(res.Data), resp.WithMeta(res.Pagination))
}
