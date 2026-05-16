package repository

import (
	"errors"

	customermodel "tech-challenge-users/internal/adapter/database/model/customer"
	"tech-challenge-users/internal/application/ports"
	"tech-challenge-users/internal/domain"
	"tech-challenge-users/pkg/converters"

	"gorm.io/gorm"
)

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) ports.CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) Create(customer *domain.Customer) error {
	m := converters.CustomerToModel(*customer)
	if err := r.db.Create(&m).Error; err != nil {
		return err
	}
	customer.ID = m.ID
	customer.CreatedAt = m.CreatedAt
	customer.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *customerRepository) FindByID(id int64) (*domain.Customer, error) {
	var m customermodel.Model
	if err := r.db.First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	c := converters.CustomerToDomain(m)
	return &c, nil
}

func (r *customerRepository) FindAll(filters ports.CustomerFilters) ([]domain.Customer, error) {
	var models []customermodel.Model

	query := r.db.Model(&customermodel.Model{}).
		Joins(`JOIN "Person" ON "Person".id = "Customer".person_id AND "Person".deleted_at IS NULL`)

	if filters.Document != "" {
		query = query.Where(`"Person".document = ?`, filters.Document)
	}
	if filters.Name != "" {
		query = query.Where(`"Person".name ILIKE ?`, "%"+filters.Name+"%")
	}
	if filters.Email != "" {
		query = query.Where(`"Person".email ILIKE ?`, "%"+filters.Email+"%")
	}
	if filters.IsActive != nil {
		query = query.Where(`"Person".is_active = ?`, *filters.IsActive)
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	customers := make([]domain.Customer, len(models))
	for i, m := range models {
		customers[i] = converters.CustomerToDomain(m)
	}
	return customers, nil
}

func (r *customerRepository) Update(customer *domain.Customer) error {
	m := converters.CustomerToModel(*customer)
	if err := r.db.Save(&m).Error; err != nil {
		return err
	}
	customer.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *customerRepository) Delete(id int64) error {
	return r.db.Delete(&customermodel.Model{}, id).Error
}
