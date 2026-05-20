package domain

import "time"

type Company struct {
	ID        int64
	Name      string
	Email     string
	Document  string
	Contact   string
	Address   Address
	CreatedAt time.Time
	UpdatedAt time.Time
}
