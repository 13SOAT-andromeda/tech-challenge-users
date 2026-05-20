package domain

import "time"

type Vehicle struct {
	ID        int64
	Plate     string
	Name      string
	Year      int
	Brand     string
	Color     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
