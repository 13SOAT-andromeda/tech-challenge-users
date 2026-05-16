package services

import (
	"errors"

	"tech-challenge-users/internal/application/ports"
	"tech-challenge-users/internal/domain"
)

var (
	ErrVehicleNotFound    = errors.New("vehicle not found")
	ErrPlateAlreadyExists = errors.New("plate already in use")
)

type CreateVehicleInput struct {
	Plate string
	Name  string
	Year  int
	Brand string
	Color string
}

type UpdateVehicleInput struct {
	Plate *string
	Name  *string
	Year  *int
	Brand *string
	Color *string
}

type VehicleService struct {
	vehicleRepo ports.VehicleRepository
}

func NewVehicleService(vehicleRepo ports.VehicleRepository) *VehicleService {
	return &VehicleService{vehicleRepo: vehicleRepo}
}

func (s *VehicleService) CreateVehicle(input CreateVehicleInput) (*domain.Vehicle, error) {
	plate, err := domain.NewPlate(input.Plate)
	if err != nil {
		return nil, err
	}

	existing, err := s.vehicleRepo.FindByPlate(plate.String())
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrPlateAlreadyExists
	}

	vehicle := domain.Vehicle{
		Plate: plate.String(),
		Name:  input.Name,
		Year:  input.Year,
		Brand: input.Brand,
		Color: input.Color,
	}
	if err := s.vehicleRepo.Create(&vehicle); err != nil {
		return nil, err
	}
	return &vehicle, nil
}

func (s *VehicleService) GetByID(id int64) (*domain.Vehicle, error) {
	vehicle, err := s.vehicleRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if vehicle == nil {
		return nil, ErrVehicleNotFound
	}
	return vehicle, nil
}

func (s *VehicleService) ListVehicles(filters ports.VehicleFilters) ([]domain.Vehicle, error) {
	return s.vehicleRepo.FindAll(filters)
}

func (s *VehicleService) UpdateVehicle(id int64, input UpdateVehicleInput) (*domain.Vehicle, error) {
	vehicle, err := s.vehicleRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if vehicle == nil {
		return nil, ErrVehicleNotFound
	}

	if input.Plate != nil {
		plate, err := domain.NewPlate(*input.Plate)
		if err != nil {
			return nil, err
		}
		normalized := plate.String()
		if normalized != vehicle.Plate {
			existing, err := s.vehicleRepo.FindByPlate(normalized)
			if err != nil {
				return nil, err
			}
			if existing != nil {
				return nil, ErrPlateAlreadyExists
			}
		}
		vehicle.Plate = normalized
	}
	if input.Name != nil {
		vehicle.Name = *input.Name
	}
	if input.Year != nil {
		vehicle.Year = *input.Year
	}
	if input.Brand != nil {
		vehicle.Brand = *input.Brand
	}
	if input.Color != nil {
		vehicle.Color = *input.Color
	}

	if err := s.vehicleRepo.Update(vehicle); err != nil {
		return nil, err
	}
	return vehicle, nil
}

func (s *VehicleService) DeleteVehicle(id int64) error {
	vehicle, err := s.vehicleRepo.FindByID(id)
	if err != nil {
		return err
	}
	if vehicle == nil {
		return ErrVehicleNotFound
	}
	return s.vehicleRepo.Delete(id)
}
