package response

import "github.com/gofiber/fiber/v2"

// SuccessOption an option for success response.
type SuccessOption func(*appSuccess)

// appSuccess standard success response that should be used in every response
// for all handlers.
type appSuccess struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Meta    any    `json:"meta,omitempty"`
}

// Success return json response with standard success response as the
// structure.
func Success(c *fiber.Ctx, options ...SuccessOption) error {
	appScs := new(appSuccess)
	appScs.Status = "SUCCESS"

	// apply all available options
	for _, opt := range options {
		opt(appScs)
	}

	return c.Status(fiber.StatusOK).JSON(*appScs)
}
