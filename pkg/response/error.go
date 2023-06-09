package response

import "github.com/gofiber/fiber/v2"

// AppErrorOption an option for error response.
type AppErrorOption func(*appError)

// appError standard error response that should be used in every response for
// all handlers.
type appError struct {
	Status  string `json:"status"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Detail  any    `json:"detail,omitempty"`
}

// Error return json response with standard error response as the structure.
func Error(c *fiber.Ctx, options ...AppErrorOption) error {
	err := new(appError)
	err.Status = "FAILED"

	// apply all available options
	for _, opt := range options {
		opt(err)
	}

	return c.Status(fiber.StatusBadRequest).JSON(*err)
}

// ErrorCode return json response with standard error response as the structure
// and additionally set the response status code.
func ErrorCode(c *fiber.Ctx, code int, options ...AppErrorOption) error {
	err := new(appError)
	err.Status = "FAILED"

	// apply all available options
	for _, opt := range options {
		opt(err)
	}

	return c.Status(code).JSON(*err)
}
