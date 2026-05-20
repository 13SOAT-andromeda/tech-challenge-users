package repository

import (
	"errors"

	employeemodel "tech-challenge-users/internal/adapter/database/model/employee"
	"tech-challenge-users/internal/application/ports"
	"tech-challenge-users/internal/domain"
	"tech-challenge-users/pkg/converters"

	"gorm.io/gorm"
)

type employeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) ports.EmployeeRepository {
	return &employeeRepository{db: db}
}

func (r *employeeRepository) Create(employee *domain.Employee) error {
	m := converters.EmployeeToModel(*employee)
	if err := r.db.Create(&m).Error; err != nil {
		return err
	}
	employee.ID = m.ID
	employee.CreatedAt = m.CreatedAt
	employee.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *employeeRepository) FindByID(id int64) (*domain.Employee, error) {
	var m employeemodel.Model
	if err := r.db.First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	e := converters.EmployeeToDomain(m)
	return &e, nil
}

func (r *employeeRepository) FindByPersonID(personID int64) (*domain.Employee, error) {
	var m employeemodel.Model
	err := r.db.Where("person_id = ?", personID).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	e := converters.EmployeeToDomain(m)
	return &e, nil
}
