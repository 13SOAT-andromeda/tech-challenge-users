package converters

import (
	vehiclemodel "tech-challenge-users/internal/adapter/database/model/vehicle"
	"tech-challenge-users/internal/domain"
)

func VehicleToModel(v domain.Vehicle) vehiclemodel.Model {
	return vehiclemodel.Model{
		ID:    v.ID,
		Plate: v.Plate,
		Name:  v.Name,
		Year:  v.Year,
		Brand: v.Brand,
		Color: v.Color,
	}
}

func VehicleToDomain(m vehiclemodel.Model) domain.Vehicle {
	return domain.Vehicle{
		ID:        m.ID,
		Plate:     m.Plate,
		Name:      m.Name,
		Year:      m.Year,
		Brand:     m.Brand,
		Color:     m.Color,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
