package ports

import "tech-challenge-users/internal/domain"

type UserFilters struct {
	Name    string
	Email   string
	Contact string
}

type UserRepository interface {
	Create(user *domain.User) error
	FindByID(id int64) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
	GetByPersonID(personID int64) (*domain.User, error)
	FindAll(filters UserFilters) ([]domain.User, error)
	Update(user *domain.User) error
	Delete(id int64) error
}
