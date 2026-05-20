package repository

import (
	"errors"

	personmodel "tech-challenge-users/internal/adapter/database/model/person"
	"tech-challenge-users/internal/application/ports"
	"tech-challenge-users/internal/domain"
	"tech-challenge-users/pkg/converters"

	"gorm.io/gorm"
)

type personRepository struct {
	db *gorm.DB
}

func NewPersonRepository(db *gorm.DB) ports.PersonRepository {
	return &personRepository{db: db}
}

func (r *personRepository) Create(person *domain.Person) error {
	m := converters.PersonToModel(*person)
	if err := r.db.Create(&m).Error; err != nil {
		return err
	}
	person.ID = m.ID
	person.CreatedAt = m.CreatedAt
	person.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *personRepository) FindByID(id int64) (*domain.Person, error) {
	var m personmodel.Model
	if err := r.db.First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	p := converters.PersonToDomain(m)
	return &p, nil
}

func (r *personRepository) GetByEmail(email string) (*domain.Person, error) {
	var m personmodel.Model
	err := r.db.Where("email = ?", email).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	p := converters.PersonToDomain(m)
	return &p, nil
}

func (r *personRepository) GetByDocument(document string) (*domain.Person, error) {
	var m personmodel.Model
	err := r.db.Where("document = ?", document).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	p := converters.PersonToDomain(m)
	return &p, nil
}

func (r *personRepository) Update(person *domain.Person) error {
	m := converters.PersonToModel(*person)
	if err := r.db.Save(&m).Error; err != nil {
		return err
	}
	person.UpdatedAt = m.UpdatedAt
	return nil
}
