package http

import "time"

type vehicleRequest struct {
	Plate        string `json:"plate"`
	InternalCode string `json:"internal_code"`
	VehicleType  string `json:"vehicle_type"`
	Brand        string `json:"brand"`
	Model        string `json:"model"`
	Capacity     int    `json:"capacity"`
	Status       string `json:"status"`
}

type vehicleResponse struct {
	ID           string    `json:"id"`
	CompanyID    string    `json:"company_id"`
	Plate        string    `json:"plate"`
	InternalCode string    `json:"internal_code"`
	VehicleType  string    `json:"vehicle_type"`
	Brand        string    `json:"brand"`
	Model        string    `json:"model"`
	Capacity     int       `json:"capacity"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type vehicleEnvelope struct {
	Vehicle vehicleResponse `json:"vehicle"`
}

type vehiclesEnvelope struct {
	Vehicles []vehicleResponse `json:"vehicles"`
}

type errorResponse struct {
	Error string `json:"error"`
}
