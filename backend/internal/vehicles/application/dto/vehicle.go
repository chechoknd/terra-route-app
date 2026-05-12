package dto

import "time"

type Vehicle struct {
	ID           string
	CompanyID    string
	Plate        string
	InternalCode string
	VehicleType  string
	Brand        string
	Model        string
	Capacity     int
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type CreateVehicleRequest struct {
	CompanyID    string
	Plate        string
	InternalCode string
	VehicleType  string
	Brand        string
	Model        string
	Capacity     int
	Status       string
}

type UpdateVehicleRequest struct {
	CompanyID    string
	ID           string
	Plate        string
	InternalCode string
	VehicleType  string
	Brand        string
	Model        string
	Capacity     int
	Status       string
}

type GetVehicleRequest struct {
	CompanyID string
	ID        string
}

type ListVehiclesRequest struct {
	CompanyID string
}

type DeactivateVehicleRequest struct {
	CompanyID string
	ID        string
}
