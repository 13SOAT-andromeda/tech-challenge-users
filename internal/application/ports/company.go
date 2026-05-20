package ports

import "tech-challenge-users/internal/domain"

type CompanyRepository interface {
	Create(company *domain.Company) error
	FindByID(id int64) (*domain.Company, error)
	FindByDocument(document string) (*domain.Company, error)
	Update(company *domain.Company) error
	Delete(id int64) error
}
