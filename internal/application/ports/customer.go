package ports

import "tech-challenge-users/internal/domain"

type CustomerFilters struct {
	Document string
	Name     string
	Email    string
	IsActive *bool
}

type CustomerRepository interface {
	Create(customer *domain.Customer) error
	FindByID(id int64) (*domain.Customer, error)
	FindAll(filters CustomerFilters) ([]domain.Customer, error)
	Update(customer *domain.Customer) error
	Delete(id int64) error
}
