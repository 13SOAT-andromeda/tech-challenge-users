package usecases

import (
	"errors"

	"tech-challenge-users/internal/application/ports"
	"tech-challenge-users/internal/domain"
)

var (
	ErrCustomerNotFound         = errors.New("customer not found")
	ErrVehicleNotFound          = errors.New("vehicle not found")
	ErrAssociationAlreadyExists = errors.New("vehicle already associated with this customer")
	ErrAssociationNotFound      = errors.New("association not found")
)

type CustomerUseCase struct {
	customerRepo ports.CustomerRepository
	vehicleRepo  ports.VehicleRepository
	cvRepo       ports.CustomerVehicleRepository
}

func NewCustomerUseCase(
	customerRepo ports.CustomerRepository,
	vehicleRepo ports.VehicleRepository,
	cvRepo ports.CustomerVehicleRepository,
) *CustomerUseCase {
	return &CustomerUseCase{
		customerRepo: customerRepo,
		vehicleRepo:  vehicleRepo,
		cvRepo:       cvRepo,
	}
}

func (uc *CustomerUseCase) AddVehicle(customerID, vehicleID int64) (*domain.CustomerVehicle, error) {
	customer, err := uc.customerRepo.FindByID(customerID)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, ErrCustomerNotFound
	}

	vehicle, err := uc.vehicleRepo.FindByID(vehicleID)
	if err != nil {
		return nil, err
	}
	if vehicle == nil {
		return nil, ErrVehicleNotFound
	}

	existing, err := uc.cvRepo.FindByCustomerAndVehicle(customerID, vehicleID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrAssociationAlreadyExists
	}

	cv := domain.CustomerVehicle{
		CustomerID: customerID,
		VehicleID:  vehicleID,
		Vehicle:    *vehicle,
	}
	if err := uc.cvRepo.Associate(&cv); err != nil {
		return nil, err
	}
	return &cv, nil
}

func (uc *CustomerUseCase) RemoveVehicle(customerID, vehicleID int64) error {
	existing, err := uc.cvRepo.FindByCustomerAndVehicle(customerID, vehicleID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrAssociationNotFound
	}
	return uc.cvRepo.Dissociate(customerID, vehicleID)
}

func (uc *CustomerUseCase) GetCustomerVehicles(customerID int64) ([]domain.CustomerVehicle, error) {
	customer, err := uc.customerRepo.FindByID(customerID)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, ErrCustomerNotFound
	}
	return uc.cvRepo.FindByCustomerID(customerID)
}
