package password

import (
	"context"

	"github.com/mdanialr/pwman_backend/internal/entity"
	repo "github.com/mdanialr/pwman_backend/internal/repository"
)

// Repository signature that's used in password domain for repository layer.
type Repository interface {
	// GetCategoryByID retrieve an entity.Category by given id.
	GetCategoryByID(ctx context.Context, id uint, opts ...repo.Options) (*entity.Category, error)
	// FindCategories retrieve all entity.Category that match given condition
	// in opts.
	FindCategories(ctx context.Context, opts ...repo.Options) ([]*entity.Category, error)
	// CreateCategory create new entity.Category and return the newly created
	// object along with assigned id as primary key.
	CreateCategory(ctx context.Context, obj entity.Category) (*entity.Category, error)
}
