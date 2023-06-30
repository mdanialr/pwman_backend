package middleware

import (
	resp "github.com/mdanialr/pwman_backend/pkg/response"

	"github.com/gofiber/fiber/v2"
	jwtMiddleware "github.com/gofiber/jwt/v3"
	"github.com/spf13/viper"
)

const (
	InvalidToken = "Invalid or Expired Token"
)

// JWT middleware that use JSON Web Token as access token.
func JWT(v *viper.Viper) fiber.Handler {
	return jwtMiddleware.New(jwtMiddleware.Config{
		ContextKey:    "jwt",
		SigningMethod: "HS256",
		SigningKey:    []byte(v.GetString("jwt.secret")),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return resp.ErrorCode(c, fiber.StatusUnauthorized, resp.WithErrMsg(InvalidToken))
		},
	})
}
