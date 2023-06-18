package password

import (
	"mime/multipart"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Request standard request object that may be used in password domain.
type Request struct {
	pagination
}

// RequestCategory standard request object that may be used in password domain.
type RequestCategory struct {
	pagination
	// ID unique identifier of each Category. Should be required when updating.
	ID uint `form:"id"`
	// Name the name of category.
	Name string `form:"name" validate:"required"`
	// Image binary file for image field that should be parsed manually from
	// delivery.
	Image *multipart.FileHeader `form:"image"`
	// Icon binary file for icon field that should be parsed manually from
	// delivery.
	Icon *multipart.FileHeader `form:"icon"`
}

// Validate apply validation rules for RequestCategory.
func (r *RequestCategory) Validate() validator.ValidationErrors {
	v := validator.New()
	v.RegisterStructValidation(r.imageValidation, RequestCategory{})

	if err := v.Struct(r); err != nil {
		return err.(validator.ValidationErrors)
	}
	return nil
}

// ValidateCreate apply validation rules for RequestCategory in create
// endpoint.
func (r *RequestCategory) ValidateCreate() validator.ValidationErrors {
	v := validator.New()
	v.RegisterStructValidation(r.createRequiredValidation, RequestCategory{})
	if err := v.Struct(r); err != nil {
		return err.(validator.ValidationErrors)
	}

	// then add the basic validation
	return r.Validate()
}

// ValidateUpdate apply validation rules for RequestCategory in update
// endpoint.
func (r *RequestCategory) ValidateUpdate() validator.ValidationErrors {
	v := validator.New()
	v.RegisterStructValidation(r.updateRequiredValidation, RequestCategory{})
	if err := v.Struct(r); err != nil {
		return err.(validator.ValidationErrors)
	}

	// then add the basic validation
	return r.Validate()
}

// NormalizeName transform value of Name field to upper-cased.
func (r *RequestCategory) NormalizeName() {
	r.Name = strings.ToUpper(r.Name)
}

// acceptedImages list of accepted image content-type.
var acceptedImages = map[string]any{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
}

// imageValidation custom validation to make sure valid image extension are
// sent.
func (r *RequestCategory) imageValidation(sl validator.StructLevel) {
	req := sl.Current().Interface().(RequestCategory)

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

// createRequiredValidation custom required fields validation in create
// endpoint.
func (r *RequestCategory) createRequiredValidation(sl validator.StructLevel) {
	req := sl.Current().Interface().(RequestCategory)

	// required for field Image
	if req.Image == nil {
		sl.ReportError(req.Image, "image", "Image", "required", "Image")
	}
	// required for field Icon
	if req.Icon == nil {
		sl.ReportError(req.Icon, "icon", "Icon", "required", "Icon")
	}
}

// updateRequiredValidation custom required fields validation in update
// endpoint.
func (r *RequestCategory) updateRequiredValidation(sl validator.StructLevel) {
	req := sl.Current().Interface().(RequestCategory)

	// required for field ID
	if req.ID < 1 {
		sl.ReportError(req.ID, "id", "ID", "required", "ID")
	}
}
