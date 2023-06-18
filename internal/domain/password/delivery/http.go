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

	apiCat := app.Group("/category", md.JWT(conf))
	apiCat.Get("/", d.IndexCategory)
	apiCat.Post("/create", d.CreateCategory)
	apiCat.Post("/update", d.UpdateCategory)
	apiCat.Post("/delete", d.DeleteCategory)

	api := app.Group("/password", md.JWT(conf))
	api.Get("/", d.Index)
	api.Post("/create", d.Create)
	api.Post("/update", d.Update)
	api.Post("/delete", d.Delete)
}

type delivery struct {
	uc pwUC.UseCase
}

func (d *delivery) Index(c *fiber.Ctx) error {
	var req pw.Request
	c.QueryParser(&req)
	// set up the query order and sort
	req.SetQuery()

	res, err := d.uc.IndexPassword(c.Context(), req)
	if err != nil {
		return resp.Error(c, resp.WithErr(err))
	}

	return resp.Success(c, resp.WithData(res.Data), resp.WithMeta(res.Pagination))
}

func (d *delivery) Create(c *fiber.Ctx) error {
	var req pw.Request
	c.BodyParser(&req)

	// validate the request
	if err := req.Validate(); err != nil {
		return resp.Error(c, resp.WithErrValidation(err))
	}

	res, err := d.uc.SavePassword(c.Context(), req)
	if err != nil {
		return resp.Error(c, resp.WithErr(err))
	}

	return resp.Success(c, resp.WithData(res))
}

func (d *delivery) Update(c *fiber.Ctx) error {
	var req pw.Request
	c.BodyParser(&req)

	// validate the request
	if err := req.ValidateUpdate(); err != nil {
		return resp.Error(c, resp.WithErrValidation(err))
	}

	if err := d.uc.UpdatePassword(c.Context(), req.ID, req); err != nil {
		return resp.Error(c, resp.WithErr(err))
	}

	return resp.Success(c, resp.WithMsg("updated successfully"))
}

func (d *delivery) Delete(c *fiber.Ctx) error {
	var req pw.Request
	c.BodyParser(&req)

	// validate the request
	if err := req.ValidateDelete(); err != nil {
		return resp.Error(c, resp.WithErrValidation(err))
	}

	if err := d.uc.DeletePassword(c.Context(), req.ID); err != nil {
		return resp.Error(c, resp.WithErr(err))
	}

	return resp.Success(c, resp.WithMsg("deleted successfully"))
}

func (d *delivery) IndexCategory(c *fiber.Ctx) error {
	var req pw.RequestCategory
	c.QueryParser(&req)
	// set up the query order and sort
	req.SetQuery()

	res, err := d.uc.IndexCategory(c.Context(), req)
	if err != nil {
		return resp.Error(c, resp.WithErr(err))
	}

	return resp.Success(c, resp.WithData(res.Data), resp.WithMeta(res.Pagination))
}

func (d *delivery) CreateCategory(c *fiber.Ctx) error {
	var req pw.RequestCategory
	c.BodyParser(&req)
	// manually retrieve binary files
	req.Icon, _ = c.FormFile("icon")
	req.Image, _ = c.FormFile("image")

	// validate the request
	if err := req.ValidateCreate(); err != nil {
		return resp.Error(c, resp.WithErrValidation(err))
	}
	// normalize name field
	req.NormalizeName()

	res, err := d.uc.SaveCategory(c.Context(), req)
	if err != nil {
		return resp.Error(c, resp.WithErr(err))
	}

	return resp.Success(c, resp.WithData(res))
}

func (d *delivery) UpdateCategory(c *fiber.Ctx) error {
	var req pw.RequestCategory
	c.BodyParser(&req)
	// manually retrieve binary files
	req.Icon, _ = c.FormFile("icon")
	req.Image, _ = c.FormFile("image")

	// validate the request
	if err := req.ValidateUpdate(); err != nil {
		return resp.Error(c, resp.WithErrValidation(err))
	}
	// normalize name field
	req.NormalizeName()

	if err := d.uc.UpdateCategory(c.Context(), req.ID, req); err != nil {
		return resp.Error(c, resp.WithErr(err))
	}

	return resp.Success(c, resp.WithMsg("updated successfully"))
}

func (d *delivery) DeleteCategory(c *fiber.Ctx) error {
	var req pw.RequestCategory
	c.BodyParser(&req)

	// validate the request
	if err := req.ValidateDelete(); err != nil {
		return resp.Error(c, resp.WithErrValidation(err))
	}

	if err := d.uc.DeleteCategory(c.Context(), req.ID); err != nil {
		return resp.Error(c, resp.WithErr(err))
	}

	return resp.Success(c, resp.WithMsg("deleted successfully"))
}
