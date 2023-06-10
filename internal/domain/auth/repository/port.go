package auth

import (
	"context"

	"github.com/mdanialr/pwman_backend/internal/entity"
)

// Repository signature that's used in auth domain for repository layer.
type Repository interface {
	// GetByCode retrieve an entity.RegisteredOTP by given code, also return
	// error if any including record not found.
	GetByCode(ctx context.Context, code string) (*entity.RegisteredOTP, error)
	// Create save new instance of entity.RegisteredOTP that's only need given
	// code.
	Create(ctx context.Context, code string) (*entity.RegisteredOTP, error)
	// DeleteAll batch delete all records of entity.RegisteredOTP.
	DeleteAll(ctx context.Context) error
}
