package auth

import (
	"context"

	"github.com/mdanialr/pwman_backend/internal/entity"

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

func (r *repository) GetByCode(ctx context.Context, code string) (*entity.RegisteredOTP, error) {
	ro := entity.RegisteredOTP{Code: code}
	return &ro, r.db.WithContext(ctx).Where(&ro).Select("id").First(&ro).Error
}

func (r *repository) Create(ctx context.Context, code string) (*entity.RegisteredOTP, error) {
	ro := entity.RegisteredOTP{Code: code}
	return &ro, r.db.WithContext(ctx).Create(&ro).Error
}

func (r *repository) DeleteAll(ctx context.Context) error {
	return r.db.WithContext(ctx).Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&entity.RegisteredOTP{}).Error
}
