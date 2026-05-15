package domain

import "time"

type CustomerVehicle struct {
	ID         int64
	CustomerID int64
	VehicleID  int64
	Vehicle    Vehicle
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
