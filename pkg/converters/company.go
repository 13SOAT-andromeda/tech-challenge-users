package converters

import (
	companymodel "tech-challenge-users/internal/adapter/database/model/company"
	"tech-challenge-users/internal/domain"
)

func CompanyToModel(c domain.Company) companymodel.Model {
	m := companymodel.Model{
		ID:       c.ID,
		Name:     c.Name,
		Email:    c.Email,
		Document: c.Document,
		Contact:  c.Contact,
	}
	if c.Address.Street != "" {
		s := c.Address.Street
		m.Address = &s
	}
	if c.Address.Number != "" {
		s := c.Address.Number
		m.AddressNumber = &s
	}
	if c.Address.Neighborhood != "" {
		s := c.Address.Neighborhood
		m.Neighborhood = &s
	}
	if c.Address.City != "" {
		s := c.Address.City
		m.City = &s
	}
	if c.Address.Country != "" {
		s := c.Address.Country
		m.Country = &s
	}
	if c.Address.ZipCode != "" {
		s := c.Address.ZipCode
		m.ZipCode = &s
	}
	return m
}

func CompanyToDomain(m companymodel.Model) domain.Company {
	c := domain.Company{
		ID:        m.ID,
		Name:      m.Name,
		Email:     m.Email,
		Document:  m.Document,
		Contact:   m.Contact,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
	if m.Address != nil {
		c.Address.Street = *m.Address
	}
	if m.AddressNumber != nil {
		c.Address.Number = *m.AddressNumber
	}
	if m.Neighborhood != nil {
		c.Address.Neighborhood = *m.Neighborhood
	}
	if m.City != nil {
		c.Address.City = *m.City
	}
	if m.Country != nil {
		c.Address.Country = *m.Country
	}
	if m.ZipCode != nil {
		c.Address.ZipCode = *m.ZipCode
	}
	return c
}
