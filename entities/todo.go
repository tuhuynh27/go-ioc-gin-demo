package entities

import (
	"time"

	"gorm.io/gorm"
)

// Todo represents the todo table structure
type Todo struct {
	gorm.Model
	ID        int       `gorm:"primaryKey"`
	Title     string    `gorm:"not null"`
	Completed bool      `gorm:"default:false"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
