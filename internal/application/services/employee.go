package services

import (
	"errors"

	"tech-challenge-users/internal/application/ports"
	"tech-challenge-users/internal/domain"
)

var ErrEmployeeNotFound = errors.New("employee not found")

type EmployeeService struct {
	employeeRepo ports.EmployeeRepository
}

func NewEmployeeService(employeeRepo ports.EmployeeRepository) *EmployeeService {
	return &EmployeeService{employeeRepo: employeeRepo}
}

func (s *EmployeeService) GetByID(id int64) (*domain.Employee, error) {
	employee, err := s.employeeRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if employee == nil {
		return nil, ErrEmployeeNotFound
	}
	return employee, nil
}
