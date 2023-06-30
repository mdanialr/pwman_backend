package seeder

import (
	"github.com/mdanialr/pwman_backend/internal/entity"

	"gorm.io/gorm"
)

func category(db *gorm.DB) {
	samples := []entity.Category{
		{
			Name:      "FAKE",
			IconPath:  "fake-icon.png",
			ImagePath: "fake.png",
		},
		{
			Name:      "DUMMIES",
			IconPath:  "dummy-icon.png",
			ImagePath: "dummy.png",
		},
	}

	for _, sample := range samples {
		db.Create(&sample)
	}
}
