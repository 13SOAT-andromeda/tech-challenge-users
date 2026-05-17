package ports

import "tech-challenge-users/internal/domain"

type EmployeeRepository interface {
	Create(employee *domain.Employee) error
	FindByID(id int64) (*domain.Employee, error)
	FindByPersonID(personID int64) (*domain.Employee, error)
}
