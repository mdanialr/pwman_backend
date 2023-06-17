package password

import (
	"mime/multipart"
	"strings"

	paginate "github.com/mdanialr/pwman_backend/pkg/pagination"

	"github.com/go-playground/validator/v10"
)

// Request standard request object that may be used in password domain.
type Request struct {
	// Name the name of category.
	Name string `form:"name" validate:"required"`
	// Image binary file for image field that should be parsed manually from
	// delivery.
	Image *multipart.FileHeader `form:"image" validate:"required"`
	// Icon binary file for icon field that should be parsed manually from
	// delivery.
	Icon *multipart.FileHeader `form:"icon" validate:"required"`
	paginate.M
	// Order the field name to query Order. Default to id.
	Order string `json:"-" query:"order"`
	// Sort to query Order. Should be filled with either asc or desc. Default
	// to asc.
	Sort string `json:"-" query:"sort"`
	// Search do search for category name from given string.
	Search string `json:"-" query:"search"`
}

// SetQuery do setup Order and Sort.
func (r *Request) SetQuery() {
	if r.Order == "" {
		r.Order = "id" // set default to id
	}
	// sanitize Sort
	r.Sort = r.sanitizeQuerySort()
	if r.Sort == "" {
		r.Sort = "asc" // set default to asc
	}
	// make sure the Sort is upper-cased
	r.Sort = strings.ToUpper(r.Sort)
}

// sanitizeQuerySort make sure Sort has the expected value.
func (r *Request) sanitizeQuerySort() string {
	switch strings.ToLower(r.Sort) {
	case "asc", "desc":
		return r.Sort
	}
	return ""
}

// Validate apply validation rules for Request.
func (r *Request) Validate() validator.ValidationErrors {
	v := validator.New()
	v.RegisterStructValidation(ImageValidation, Request{})

	if err := v.Struct(r); err != nil {
		return err.(validator.ValidationErrors)
	}
	return nil
}

// NormalizeName transform value of Name field to upper-cased.
func (r *Request) NormalizeName() {
	r.Name = strings.ToUpper(r.Name)
}

// acceptedImages list of accepted image content-type.
var acceptedImages = map[string]any{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
}

// ImageValidation custom validation to make sure valid image extension are
// sent.
func ImageValidation(sl validator.StructLevel) {
	req := sl.Current().Interface().(Request)

	if req.Image != nil {
		ext := req.Image.Header.Get("content-type")
		if _, ok := acceptedImages[ext]; !ok {
			sl.ReportError(req.Image, "image", "Image", "image", "Image")
		}
	}
	if req.Icon != nil {
		ext := req.Icon.Header.Get("content-type")
		if _, ok := acceptedImages[ext]; !ok {
			sl.ReportError(req.Icon, "icon", "Icon", "image", "Icon")
		}
	}
}
