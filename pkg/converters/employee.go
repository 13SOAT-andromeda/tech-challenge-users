package converters

import (
	employeemodel "tech-challenge-users/internal/adapter/database/model/employee"
	"tech-challenge-users/internal/domain"
)

func EmployeeToModel(e domain.Employee) employeemodel.Model {
	return employeemodel.Model{
		ID:       e.ID,
		Position: e.Position,
		PersonID: e.PersonID,
	}
}

func EmployeeToDomain(m employeemodel.Model) domain.Employee {
	return domain.Employee{
		ID:        m.ID,
		Position:  m.Position,
		PersonID:  m.PersonID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
