package domain

import (
	"errors"
	"time"
)

const (
	RoleCustomer      = "customer"
	RoleAttendant     = "attendant"
	RoleMechanic      = "mechanic"
	RoleAdministrator = "administrator"
)

var validRoles = map[string]struct{}{
	RoleCustomer:      {},
	RoleAttendant:     {},
	RoleMechanic:      {},
	RoleAdministrator: {},
}

type User struct {
	ID        int64
	Password  string
	Role      string
	PersonID  int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func ValidateRole(role string) error {
	if _, ok := validRoles[role]; !ok {
		return errors.New("role must be one of: customer, attendant, mechanic, administrator")
	}
	return nil
}
