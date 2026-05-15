package company

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID           int64          `gorm:"primaryKey;autoIncrement"`
	Name         string         `gorm:"not null"`
	Email        string         `gorm:"not null"`
	Document     string         `gorm:"not null"`
	Contact      string         `gorm:"not null"`
	Address      *string
	AddressNumber *string
	Neighborhood *string
	City         *string
	Country      *string
	ZipCode      *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (Model) TableName() string { return "Company" }
