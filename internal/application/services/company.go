package services

import (
	"errors"

	"tech-challenge-users/internal/application/ports"
	"tech-challenge-users/internal/domain"
)

var ErrCompanyNotFound = errors.New("company not found")

type CreateCompanyInput struct {
	Name     string
	Email    string
	Document string
	Contact  string
	Address  domain.Address
}

type UpdateCompanyInput struct {
	Name    *string
	Email   *string
	Contact *string
	Address *domain.Address
}

type CompanyService struct {
	companyRepo ports.CompanyRepository
}

func NewCompanyService(companyRepo ports.CompanyRepository) *CompanyService {
	return &CompanyService{companyRepo: companyRepo}
}

func (s *CompanyService) CreateCompany(input CreateCompanyInput) (*domain.Company, error) {
	if _, err := domain.NewDocument(input.Document); err != nil {
		return nil, err
	}

	existing, err := s.companyRepo.FindByDocument(input.Document)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrDocumentAlreadyExists
	}

	company := domain.Company{
		Name:     input.Name,
		Email:    input.Email,
		Document: input.Document,
		Contact:  input.Contact,
		Address:  input.Address,
	}
	if err := s.companyRepo.Create(&company); err != nil {
		return nil, err
	}
	return &company, nil
}

func (s *CompanyService) GetByID(id int64) (*domain.Company, error) {
	company, err := s.companyRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if company == nil {
		return nil, ErrCompanyNotFound
	}
	return company, nil
}

func (s *CompanyService) UpdateCompany(id int64, input UpdateCompanyInput) (*domain.Company, error) {
	company, err := s.companyRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if company == nil {
		return nil, ErrCompanyNotFound
	}

	if input.Name != nil {
		company.Name = *input.Name
	}
	if input.Email != nil {
		company.Email = *input.Email
	}
	if input.Contact != nil {
		company.Contact = *input.Contact
	}
	if input.Address != nil {
		company.Address = *input.Address
	}

	if err := s.companyRepo.Update(company); err != nil {
		return nil, err
	}
	return company, nil
}

func (s *CompanyService) DeleteCompany(id int64) error {
	company, err := s.companyRepo.FindByID(id)
	if err != nil {
		return err
	}
	if company == nil {
		return ErrCompanyNotFound
	}
	return s.companyRepo.Delete(id)
}
