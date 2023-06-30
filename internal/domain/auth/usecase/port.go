package auth

import (
	"context"

	"github.com/mdanialr/pwman_backend/internal/domain/auth"
)

// UseCase a use case spec that's used in authentication domain.
type UseCase interface {
	// ValidateOTP return a Response by given request.
	ValidateOTP(ctx context.Context, req auth.Request) (*auth.Response, error)
	// CreateJWT create new jwt claims, then append the token to Response.
	CreateJWT(ctx context.Context) (*auth.Response, error)
}
