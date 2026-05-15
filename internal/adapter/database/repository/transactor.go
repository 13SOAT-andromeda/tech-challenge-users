package repository

import (
	"tech-challenge-users/internal/application/ports"

	"gorm.io/gorm"
)

type transactor struct {
	db *gorm.DB
}

func NewTransactor(db *gorm.DB) ports.Transactor {
	return &transactor{db: db}
}

func (t *transactor) WithinTransaction(fn func(repos ports.Repositories) error) error {
	return t.db.Transaction(func(tx *gorm.DB) error {
		return fn(ports.Repositories{
			Person:   NewPersonRepository(tx),
			User:     NewUserRepository(tx),
			Employee: NewEmployeeRepository(tx),
		})
	})
}
