package http

import "time"

type driverRequest struct {
	UserID         string `json:"user_id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	DocumentNumber string `json:"document_number"`
	Phone          string `json:"phone"`
	Email          string `json:"email"`
	LicenseNumber  string `json:"license_number"`
	Status         string `json:"status"`
}

type driverResponse struct {
	ID             string    `json:"id"`
	CompanyID      string    `json:"company_id"`
	UserID         string    `json:"user_id,omitempty"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	DocumentNumber string    `json:"document_number"`
	Phone          string    `json:"phone"`
	Email          string    `json:"email,omitempty"`
	LicenseNumber  string    `json:"license_number"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type driverEnvelope struct {
	Driver driverResponse `json:"driver"`
}

type driversEnvelope struct {
	Drivers []driverResponse `json:"drivers"`
}

type errorResponse struct {
	Error string `json:"error"`
}
