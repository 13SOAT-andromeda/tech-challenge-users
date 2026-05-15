package ports

import "tech-challenge-users/internal/domain"

type PersonRepository interface {
	Create(person *domain.Person) error
	FindByID(id int64) (*domain.Person, error)
	GetByEmail(email string) (*domain.Person, error)
	GetByDocument(document string) (*domain.Person, error)
	Update(person *domain.Person) error
}
