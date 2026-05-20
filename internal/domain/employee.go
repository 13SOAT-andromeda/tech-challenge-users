package domain

import "time"

type Employee struct {
	ID        int64
	Position  string
	PersonID  int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
