package converters

import (
	cvmodel "tech-challenge-users/internal/adapter/database/model/customer_vehicle"
	"tech-challenge-users/internal/domain"
)

func CustomerVehicleToModel(cv domain.CustomerVehicle) cvmodel.Model {
	return cvmodel.Model{
		ID:         cv.ID,
		CustomerID: cv.CustomerID,
		VehicleID:  cv.VehicleID,
	}
}

func CustomerVehicleToDomain(m cvmodel.Model) domain.CustomerVehicle {
	return domain.CustomerVehicle{
		ID:         m.ID,
		CustomerID: m.CustomerID,
		VehicleID:  m.VehicleID,
		Vehicle:    VehicleToDomain(m.Vehicle),
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}
