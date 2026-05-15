package domain

import (
	"errors"
	"time"
)

const (
	CustomerTypePF = "pf"
	CustomerTypePJ = "pj"
)

type Customer struct {
	ID        int64
	Type      string
	PersonID  int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func ValidateCustomerType(t string) error {
	if t != CustomerTypePF && t != CustomerTypePJ {
		return errors.New("type must be 'pf' or 'pj'")
	}
	return nil
}
