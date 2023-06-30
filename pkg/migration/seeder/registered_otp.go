package seeder

import (
	"strconv"

	"github.com/mdanialr/pwman_backend/internal/entity"

	"gorm.io/gorm"
)

func registeredOtp(db *gorm.DB) {
	samples := []uint{123456, 208765, 872532}

	for _, sample := range samples {
		ro := entity.RegisteredOTP{Code: strconv.Itoa(int(sample))}
		db.Create(&ro)
	}
}
