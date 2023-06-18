package password

import (
	"context"

	"github.com/mdanialr/pwman_backend/internal/entity"
	repo "github.com/mdanialr/pwman_backend/internal/repository"
)

// Repository signature that's used in password domain for repository layer.
type Repository interface {
	// GetPasswordByID retrieve an entity.Password by given id.
	GetPasswordByID(ctx context.Context, id uint, opts ...repo.Options) (*entity.Password, error)
	// FindPassword retrieve all entity.Password that match given condition in
	// opts.
	FindPassword(ctx context.Context, opts ...repo.Options) ([]*entity.Password, error)
	// CreatePassword create new entity.Password and return the newly created
	// object along with assigned id as primary key.
	CreatePassword(ctx context.Context, obj entity.Password) (*entity.Password, error)
	// UpdatePassword update existing entity.Password that match given id and
	// return the updated object.
	UpdatePassword(ctx context.Context, id uint, obj entity.Password, opts ...repo.Options) (*entity.Password, error)
	// DeletePassword soft delete entity.Password that match given id.
	DeletePassword(ctx context.Context, id uint) error
	// GetCategoryByID retrieve an entity.Category by given id.
	GetCategoryByID(ctx context.Context, id uint, opts ...repo.Options) (*entity.Category, error)
	// FindCategories retrieve all entity.Category that match given condition
	// in opts.
	FindCategories(ctx context.Context, opts ...repo.Options) ([]*entity.Category, error)
	// CreateCategory create new entity.Category and return the newly created
	// object along with assigned id as primary key.
	CreateCategory(ctx context.Context, obj entity.Category) (*entity.Category, error)
	// UpdateCategory update existing entity.Category that match given id and
	// return the updated object.
	UpdateCategory(ctx context.Context, id uint, obj entity.Category, opts ...repo.Options) (*entity.Category, error)
}
