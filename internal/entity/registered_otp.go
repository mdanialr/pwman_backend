package entity

import (
	"time"

	"gorm.io/gorm"
)

// RegisteredOTP object for table `registered_otp`.
type RegisteredOTP struct {
	ID        uint `gorm:"primaryKey"`
	Code      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
