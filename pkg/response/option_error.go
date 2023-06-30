package response

import (
	"github.com/go-playground/validator/v10"

	stderr "github.com/mdanialr/pwman_backend/internal/err"
)

// WithErr option to add given error message to error response as `message`
// field.
func WithErr(err error) AppErrorOption {
	return func(a *appError) {
		a.Message = err.Error()

		// check in case error is returned by use case layer
		switch e := err.(type) {
		case *stderr.UC:
			a.Code = e.Code
			a.Message = e.Msg
		}
	}
}

// WithErrMsg option to add given message to error response as `message` field.
func WithErrMsg(msg string) AppErrorOption {
	return func(a *appError) {
		a.Message = msg
	}
}

// WithErrCode option to add given code to error response as `code` field.
func WithErrCode(code string) AppErrorOption {
	return func(a *appError) {
		a.Code = code
	}
}

// WithErrDetail option to add given detail to error response as `detail`
// field.
func WithErrDetail(detail any) AppErrorOption {
	return func(a *appError) {
		a.Detail = detail
	}
}

// WithErrValidation option to add given valid to error response as `detail`
// field.
func WithErrValidation(valid validator.ValidationErrors) AppErrorOption {
	return func(a *appError) {
		a.Detail = NewValidationErrors(valid)
	}
}
