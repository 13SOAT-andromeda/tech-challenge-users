package user

import (
	"time"

	personmodel "tech-challenge-users/internal/adapter/database/model/person"

	"gorm.io/gorm"
)

type Model struct {
	ID        int64              `gorm:"primaryKey;autoIncrement"`
	Password  string             `gorm:"not null"`
	Role      string             `gorm:"not null"`
	PersonID  int64              `gorm:"not null"`
	Person    personmodel.Model  `gorm:"foreignKey:PersonID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt     `gorm:"index"`
}

func (Model) TableName() string { return "User" }
