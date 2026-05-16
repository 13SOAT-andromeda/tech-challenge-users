package domain

import "time"

type Person struct {
	ID        int64
	Name      string
	Email     string
	Contact   string
	Document  string
	IsActive  bool
	Address   Address
	CreatedAt time.Time
	UpdatedAt time.Time
}
