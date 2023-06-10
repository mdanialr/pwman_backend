package seeder

import (
	"github.com/mdanialr/pwman_backend/internal/entity"

	"gorm.io/gorm"
)

func category(db *gorm.DB) {
	samples := []entity.Category{
		{
			Name:      "Fake",
			IconPath:  "icon/fake-icon.png",
			ImagePath: "img/fake.png",
		},
		{
			Name:      "Dummies",
			IconPath:  "icon/dummy-icon.png",
			ImagePath: "img/dummy.png",
		},
	}

	for _, sample := range samples {
		db.Create(&sample)
	}
}
