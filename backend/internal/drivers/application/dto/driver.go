package dto

import "time"

type Driver struct {
	ID             string
	CompanyID      string
	UserID         string
	FirstName      string
	LastName       string
	DocumentNumber string
	Phone          string
	Email          string
	LicenseNumber  string
	Status         string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type CreateDriverRequest struct {
	CompanyID      string
	UserID         string
	FirstName      string
	LastName       string
	DocumentNumber string
	Phone          string
	Email          string
	LicenseNumber  string
	Status         string
}

type UpdateDriverRequest struct {
	CompanyID      string
	ID             string
	UserID         string
	FirstName      string
	LastName       string
	DocumentNumber string
	Phone          string
	Email          string
	LicenseNumber  string
	Status         string
}

type GetDriverRequest struct {
	CompanyID string
	ID        string
}

type ListDriversRequest struct {
	CompanyID string
}

type DeactivateDriverRequest struct {
	CompanyID string
	ID        string
}
