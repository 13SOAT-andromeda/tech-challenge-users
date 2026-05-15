package customer

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        int64          `gorm:"primaryKey;autoIncrement"`
	Type      string         `gorm:"not null"`
	PersonID  int64          `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (Model) TableName() string { return "Customer" }
