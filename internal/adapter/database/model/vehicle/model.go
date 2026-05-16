package vehicle

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        int64          `gorm:"primaryKey;autoIncrement"`
	Plate     string         `gorm:"not null"`
	Name      string         `gorm:"not null"`
	Year      int            `gorm:"not null"`
	Brand     string         `gorm:"not null"`
	Color     string         `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (Model) TableName() string { return "Vehicle" }
