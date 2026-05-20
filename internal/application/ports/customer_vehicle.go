package ports

import "tech-challenge-users/internal/domain"

type CustomerVehicleRepository interface {
	Associate(cv *domain.CustomerVehicle) error
	Dissociate(customerID, vehicleID int64) error
	FindByCustomerID(customerID int64) ([]domain.CustomerVehicle, error)
	FindByCustomerAndVehicle(customerID, vehicleID int64) (*domain.CustomerVehicle, error)
}
