package seeder

import (
	help "github.com/mdanialr/pwman_backend/pkg/helper"

	"gorm.io/gorm"
)

func users(db *gorm.DB) {
	hashedPassword, _ := help.HashPassword("password")
	userSamples := []password.Model{
		{
			Username:  "admin",
			Password:  hashedPassword,
			RoleID:    1,
			FirstName: "Robert",
			LastName:  "Amore",
			Address:   "1615 Fannie Street, San Angelo, Texas. 76903",
		},
		{
			Username:  "user",
			Password:  hashedPassword,
			RoleID:    2,
			FirstName: "David",
			LastName:  "Williams",
			Address:   "4502 Better Street, Fort Worth, Kansas",
		},
	}

	for _, usr := range userSamples {
		db.Create(&usr)
	}
}
