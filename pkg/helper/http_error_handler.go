package helper

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

// DefaultHTTPErrorHandler default HTTP error handler for fiber.Handler.
func DefaultHTTPErrorHandler(c *fiber.Ctx, err error) error {
	var e = new(fiber.Error)
	var msg string
	var code = fiber.StatusInternalServerError

	if errors.As(err, &e) {
		code = e.Code
		msg = defaultHTTPErrorHandlerMsg(code)
	}

	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	return c.Status(code).JSON(fiber.Map{"error": msg})
}

// defaultHTTPErrorHandlerMsg map given http status code to predefined message.
func defaultHTTPErrorHandlerMsg(code int) string {
	switch code {
	case fiber.StatusNotFound:
		return "Not Found"
	case fiber.StatusMethodNotAllowed:
		return "This method is not allowed here!"
	}
	return "Something was wrong!"
}
