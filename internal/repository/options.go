package repo

import (
	"github.com/mdanialr/pwman_backend/pkg/pagination"

	"gorm.io/gorm"
)

// Options signature that should be used to optionally add query to each
// repository layer implementations that interact with gorm.DB as the data
// source.
type Options func(*gorm.DB) *gorm.DB

// Cols add query Select.
// Example:
//
//	repo.Cols("id", "created_at", "updated_at")
func Cols(cols ...string) Options {
	return func(db *gorm.DB) *gorm.DB {
		return db.Select(cols)
	}
}

// Order add query Order.
//
// Example:
//
//	repo.Order("created_at DESC")
func Order(orders ...string) Options {
	return func(db *gorm.DB) *gorm.DB {
		for _, ord := range orders {
			db = db.Order(ord)
		}
		return db
	}
}

// Cons add query Where for each given cons. Each given conditions will be
// combined by GORM using AND.
//
// Example:
//
//	repo.Cons("id IS NULL"), repo.Cons("name IS NOT NULL")
func Cons(cons ...string) Options {
	return func(db *gorm.DB) *gorm.DB {
		for _, con := range cons {
			db = db.Where(con)
		}
		return db
	}
}

// EagerLoad simple preload/eager loading the given field/relation name without
// any special condition which means use the default relation as the condition.
//
// Example:
//
//	repo.EagerLoad("Role")
func EagerLoad(fields ...string) Options {
	return func(db *gorm.DB) *gorm.DB {
		for _, field := range fields {
			db = db.Preload(field)
		}
		return db
	}
}

// Paginate add query Limit & Offset accordingly by given paginate.M.
//
// Example:
//
//	repo.Paginate(&paginate.M{Limit: 4})
func Paginate(p *paginate.M) Options {
	if p == nil {
		return func(db *gorm.DB) *gorm.DB {
			return p.Set(db)
		}
	}
	return p.Set
}
