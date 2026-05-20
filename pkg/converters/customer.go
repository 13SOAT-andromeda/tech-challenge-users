package converters

import (
	customermodel "tech-challenge-users/internal/adapter/database/model/customer"
	"tech-challenge-users/internal/domain"
)

func CustomerToModel(c domain.Customer) customermodel.Model {
	return customermodel.Model{
		ID:       c.ID,
		Type:     c.Type,
		PersonID: c.PersonID,
	}
}

func CustomerToDomain(m customermodel.Model) domain.Customer {
	return domain.Customer{
		ID:        m.ID,
		Type:      m.Type,
		PersonID:  m.PersonID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
