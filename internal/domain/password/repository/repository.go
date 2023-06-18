package password

import (
	"context"

	"github.com/mdanialr/pwman_backend/internal/entity"
	repo "github.com/mdanialr/pwman_backend/internal/repository"

	"gorm.io/gorm"
)

// NewRepository return concrete implementation of Repository that use gorm.DB
// as the data source.
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

type repository struct {
	db *gorm.DB
}

func (r *repository) GetPasswordByID(ctx context.Context, id uint, opts ...repo.Options) (*entity.Password, error) {
	q := r.db.WithContext(ctx)
	p := entity.Password{ID: id}

	// apply options
	for _, opt := range opts {
		q = opt(q)
	}

	return &p, q.First(&p).Error
}

func (r *repository) FindPassword(ctx context.Context, opts ...repo.Options) ([]*entity.Password, error) {
	q := r.db.WithContext(ctx).Model(&entity.Password{})
	var p []*entity.Password

	// apply options
	for _, opt := range opts {
		q = opt(q)
	}

	return p, q.Find(&p).Error
}

func (r *repository) GetCategoryByID(ctx context.Context, id uint, opts ...repo.Options) (*entity.Category, error) {
	q := r.db.WithContext(ctx)
	c := entity.Category{ID: id}

	// apply options
	for _, opt := range opts {
		q = opt(q)
	}

	return &c, q.First(&c).Error
}

func (r *repository) FindCategories(ctx context.Context, opts ...repo.Options) ([]*entity.Category, error) {
	q := r.db.WithContext(ctx).Model(&entity.Category{})
	var c []*entity.Category

	// apply options
	for _, opt := range opts {
		q = opt(q)
	}

	return c, q.Find(&c).Error
}

func (r *repository) CreateCategory(ctx context.Context, obj entity.Category) (*entity.Category, error) {
	q := r.db.WithContext(ctx)

	return &obj, q.Create(&obj).Error
}

func (r *repository) UpdateCategory(ctx context.Context, id uint, obj entity.Category, opts ...repo.Options) (*entity.Category, error) {
	q := r.db.WithContext(ctx)
	c := entity.Category{ID: id}

	// apply options
	for _, opt := range opts {
		q = opt(q)
	}

	return &c, q.Model(&c).Updates(obj).Error
}
