package person

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID            int64          `gorm:"primaryKey;autoIncrement"`
	Name          string         `gorm:"not null"`
	Email         string         `gorm:"not null"`
	Contact       *string
	Document      string         `gorm:"not null"`
	IsActive      bool           `gorm:"not null;default:true"`
	Address       *string
	AddressNumber *string
	Neighborhood  *string
	City          *string
	Country       *string
	ZipCode       *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

func (Model) TableName() string { return "Person" }
