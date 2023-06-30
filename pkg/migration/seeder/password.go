package seeder

import (
	"github.com/mdanialr/pwman_backend/internal/entity"

	"gorm.io/gorm"
)

func password(db *gorm.DB) {
	samples := []entity.Password{
		{
			Username:   "hello-world",
			Password:   "password",
			CategoryID: 1,
		},
		{
			Username:   "hi",
			Password:   "password",
			CategoryID: 2,
		},
	}

	for _, sample := range samples {
		db.Create(&sample)
	}
}
