package entity

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID        uint `gorm:"primarykey"`
	Name      string
	ImagePath string
	IconPath  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
