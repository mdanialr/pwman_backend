package auth

import "github.com/go-playground/validator/v10"

// Request standard request object that may be used in auth domain.
type Request struct {
	Code string `json:"code" validate:"required,numeric"`
}

// Validate apply validation rules for Request.
func (r *Request) Validate() validator.ValidationErrors {
	if err := validator.New().Struct(r); err != nil {
		return err.(validator.ValidationErrors)
	}
	return nil
}
