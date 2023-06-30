package entity

import (
	"time"

	"gorm.io/gorm"
)

type Password struct {
	ID         uint `gorm:"primarykey"`
	Username   string
	Password   string
	CategoryID uint
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}
