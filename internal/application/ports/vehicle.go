package ports

import "tech-challenge-users/internal/domain"

type VehicleFilters struct {
	Name  string
	Plate string
	Brand string
	Color string
	Year  int
}

type VehicleRepository interface {
	Create(vehicle *domain.Vehicle) error
	FindByID(id int64) (*domain.Vehicle, error)
	FindByPlate(plate string) (*domain.Vehicle, error)
	FindAll(filters VehicleFilters) ([]domain.Vehicle, error)
	Update(vehicle *domain.Vehicle) error
	Delete(id int64) error
}
