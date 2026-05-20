package repository

import (
	"errors"

	vehiclemodel "tech-challenge-users/internal/adapter/database/model/vehicle"
	"tech-challenge-users/internal/application/ports"
	"tech-challenge-users/internal/domain"
	"tech-challenge-users/pkg/converters"

	"gorm.io/gorm"
)

type vehicleRepository struct {
	db *gorm.DB
}

func NewVehicleRepository(db *gorm.DB) ports.VehicleRepository {
	return &vehicleRepository{db: db}
}

func (r *vehicleRepository) Create(vehicle *domain.Vehicle) error {
	m := converters.VehicleToModel(*vehicle)
	if err := r.db.Create(&m).Error; err != nil {
		return err
	}
	vehicle.ID = m.ID
	vehicle.CreatedAt = m.CreatedAt
	vehicle.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *vehicleRepository) FindByID(id int64) (*domain.Vehicle, error) {
	var m vehiclemodel.Model
	if err := r.db.First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	v := converters.VehicleToDomain(m)
	return &v, nil
}

func (r *vehicleRepository) FindByPlate(plate string) (*domain.Vehicle, error) {
	var m vehiclemodel.Model
	if err := r.db.Where("plate = ?", plate).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	v := converters.VehicleToDomain(m)
	return &v, nil
}

func (r *vehicleRepository) FindAll(filters ports.VehicleFilters) ([]domain.Vehicle, error) {
	var models []vehiclemodel.Model
	query := r.db.Model(&vehiclemodel.Model{})

	if filters.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filters.Name+"%")
	}
	if filters.Plate != "" {
		query = query.Where("plate ILIKE ?", "%"+filters.Plate+"%")
	}
	if filters.Brand != "" {
		query = query.Where("brand ILIKE ?", "%"+filters.Brand+"%")
	}
	if filters.Color != "" {
		query = query.Where("color ILIKE ?", "%"+filters.Color+"%")
	}
	if filters.Year != 0 {
		query = query.Where("year = ?", filters.Year)
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	vehicles := make([]domain.Vehicle, len(models))
	for i, m := range models {
		vehicles[i] = converters.VehicleToDomain(m)
	}
	return vehicles, nil
}

func (r *vehicleRepository) Update(vehicle *domain.Vehicle) error {
	m := converters.VehicleToModel(*vehicle)
	if err := r.db.Save(&m).Error; err != nil {
		return err
	}
	vehicle.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *vehicleRepository) Delete(id int64) error {
	return r.db.Delete(&vehiclemodel.Model{}, id).Error
}
