package converters

import (
	personmodel "tech-challenge-users/internal/adapter/database/model/person"
	"tech-challenge-users/internal/domain"
)

func PersonToModel(p domain.Person) personmodel.Model {
	m := personmodel.Model{
		ID:       p.ID,
		Name:     p.Name,
		Email:    p.Email,
		Document: p.Document,
		IsActive: p.IsActive,
	}
	if p.Contact != "" {
		m.Contact = &p.Contact
	}
	if p.Address.Street != "" {
		s := p.Address.Street
		m.Address = &s
	}
	if p.Address.Number != "" {
		s := p.Address.Number
		m.AddressNumber = &s
	}
	if p.Address.Neighborhood != "" {
		s := p.Address.Neighborhood
		m.Neighborhood = &s
	}
	if p.Address.City != "" {
		s := p.Address.City
		m.City = &s
	}
	if p.Address.Country != "" {
		s := p.Address.Country
		m.Country = &s
	}
	if p.Address.ZipCode != "" {
		s := p.Address.ZipCode
		m.ZipCode = &s
	}
	return m
}

func PersonToDomain(m personmodel.Model) domain.Person {
	p := domain.Person{
		ID:        m.ID,
		Name:      m.Name,
		Email:     m.Email,
		Document:  m.Document,
		IsActive:  m.IsActive,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
	if m.Contact != nil {
		p.Contact = *m.Contact
	}
	if m.Address != nil {
		p.Address.Street = *m.Address
	}
	if m.AddressNumber != nil {
		p.Address.Number = *m.AddressNumber
	}
	if m.Neighborhood != nil {
		p.Address.Neighborhood = *m.Neighborhood
	}
	if m.City != nil {
		p.Address.City = *m.City
	}
	if m.Country != nil {
		p.Address.Country = *m.Country
	}
	if m.ZipCode != nil {
		p.Address.ZipCode = *m.ZipCode
	}
	return p
}
