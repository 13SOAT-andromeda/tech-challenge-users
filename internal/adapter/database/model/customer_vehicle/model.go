package customer_vehicle

import (
	"time"

	vehiclemodel "tech-challenge-users/internal/adapter/database/model/vehicle"

	"gorm.io/gorm"
)

type Model struct {
	ID         int64              `gorm:"primaryKey;autoIncrement"`
	CustomerID int64              `gorm:"not null"`
	VehicleID  int64              `gorm:"not null"`
	Vehicle    vehiclemodel.Model `gorm:"foreignKey:VehicleID"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (Model) TableName() string { return "Customer_Vehicle" }
