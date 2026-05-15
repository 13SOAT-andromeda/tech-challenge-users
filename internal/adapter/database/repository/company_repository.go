package repository

import (
	"errors"

	companymodel "tech-challenge-users/internal/adapter/database/model/company"
	"tech-challenge-users/internal/application/ports"
	"tech-challenge-users/internal/domain"
	"tech-challenge-users/pkg/converters"

	"gorm.io/gorm"
)

type companyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) ports.CompanyRepository {
	return &companyRepository{db: db}
}

func (r *companyRepository) Create(company *domain.Company) error {
	m := converters.CompanyToModel(*company)
	if err := r.db.Create(&m).Error; err != nil {
		return err
	}
	company.ID = m.ID
	company.CreatedAt = m.CreatedAt
	company.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *companyRepository) FindByID(id int64) (*domain.Company, error) {
	var m companymodel.Model
	if err := r.db.First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	c := converters.CompanyToDomain(m)
	return &c, nil
}

func (r *companyRepository) FindByDocument(document string) (*domain.Company, error) {
	var m companymodel.Model
	if err := r.db.Where("document = ?", document).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	c := converters.CompanyToDomain(m)
	return &c, nil
}

func (r *companyRepository) Update(company *domain.Company) error {
	m := converters.CompanyToModel(*company)
	if err := r.db.Save(&m).Error; err != nil {
		return err
	}
	company.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *companyRepository) Delete(id int64) error {
	return r.db.Delete(&companymodel.Model{}, id).Error
}
