package repository

import (
	"errors"

	cvmodel "tech-challenge-users/internal/adapter/database/model/customer_vehicle"
	"tech-challenge-users/internal/application/ports"
	"tech-challenge-users/internal/domain"
	"tech-challenge-users/pkg/converters"

	"gorm.io/gorm"
)

type customerVehicleRepository struct {
	db *gorm.DB
}

func NewCustomerVehicleRepository(db *gorm.DB) ports.CustomerVehicleRepository {
	return &customerVehicleRepository{db: db}
}

func (r *customerVehicleRepository) Associate(cv *domain.CustomerVehicle) error {
	m := converters.CustomerVehicleToModel(*cv)
	if err := r.db.Create(&m).Error; err != nil {
		return err
	}
	cv.ID = m.ID
	cv.CreatedAt = m.CreatedAt
	cv.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *customerVehicleRepository) Dissociate(customerID, vehicleID int64) error {
	return r.db.
		Where("customer_id = ? AND vehicle_id = ?", customerID, vehicleID).
		Delete(&cvmodel.Model{}).Error
}

func (r *customerVehicleRepository) FindByCustomerID(customerID int64) ([]domain.CustomerVehicle, error) {
	var models []cvmodel.Model
	err := r.db.Preload("Vehicle").
		Where("customer_id = ?", customerID).
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	result := make([]domain.CustomerVehicle, len(models))
	for i, m := range models {
		result[i] = converters.CustomerVehicleToDomain(m)
	}
	return result, nil
}

func (r *customerVehicleRepository) FindByCustomerAndVehicle(customerID, vehicleID int64) (*domain.CustomerVehicle, error) {
	var m cvmodel.Model
	err := r.db.
		Where("customer_id = ? AND vehicle_id = ?", customerID, vehicleID).
		First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	cv := converters.CustomerVehicleToDomain(m)
	return &cv, nil
}
