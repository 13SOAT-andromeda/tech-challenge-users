package repository

import (
	"errors"

	usermodel "tech-challenge-users/internal/adapter/database/model/user"
	"tech-challenge-users/internal/application/ports"
	"tech-challenge-users/internal/domain"
	"tech-challenge-users/pkg/converters"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) ports.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	m := converters.UserToModel(*user)
	if err := r.db.Create(&m).Error; err != nil {
		return err
	}
	user.ID = m.ID
	user.CreatedAt = m.CreatedAt
	user.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *userRepository) FindByID(id int64) (*domain.User, error) {
	var m usermodel.Model
	if err := r.db.First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	u := converters.UserToDomain(m)
	return &u, nil
}

func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	var m usermodel.Model
	err := r.db.
		Joins(`JOIN "Person" ON "Person".id = "User".person_id AND "Person".deleted_at IS NULL`).
		Where(`"Person".email = ?`, email).
		First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	u := converters.UserToDomain(m)
	return &u, nil
}

func (r *userRepository) GetByPersonID(personID int64) (*domain.User, error) {
	var m usermodel.Model
	err := r.db.Where("person_id = ?", personID).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	u := converters.UserToDomain(m)
	return &u, nil
}

func (r *userRepository) FindAll(filters ports.UserFilters) ([]domain.User, error) {
	var models []usermodel.Model
	query := r.db.Model(&usermodel.Model{})

	if filters.Name != "" || filters.Email != "" || filters.Contact != "" {
		query = query.Joins(`JOIN "Person" ON "Person".id = "User".person_id AND "Person".deleted_at IS NULL`)
		if filters.Name != "" {
			query = query.Where(`"Person".name ILIKE ?`, "%"+filters.Name+"%")
		}
		if filters.Email != "" {
			query = query.Where(`"Person".email ILIKE ?`, "%"+filters.Email+"%")
		}
		if filters.Contact != "" {
			query = query.Where(`"Person".contact ILIKE ?`, "%"+filters.Contact+"%")
		}
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	users := make([]domain.User, len(models))
	for i, m := range models {
		users[i] = converters.UserToDomain(m)
	}
	return users, nil
}

func (r *userRepository) Update(user *domain.User) error {
	m := converters.UserToModel(*user)
	if err := r.db.Save(&m).Error; err != nil {
		return err
	}
	user.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *userRepository) Delete(id int64) error {
	return r.db.Delete(&usermodel.Model{}, id).Error
}
