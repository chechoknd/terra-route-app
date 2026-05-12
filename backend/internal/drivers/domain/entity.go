package domain

import (
	"strings"
	"time"
)

type DriverStatus string

const (
	DriverStatusActive    DriverStatus = "active"
	DriverStatusInactive  DriverStatus = "inactive"
	DriverStatusSuspended DriverStatus = "suspended"
)

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
	Status         DriverStatus
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (d Driver) Validate() error {
	if strings.TrimSpace(d.CompanyID) == "" {
		return ErrInvalidDriver
	}
	if strings.TrimSpace(d.FirstName) == "" {
		return ErrInvalidDriver
	}
	if strings.TrimSpace(d.LastName) == "" {
		return ErrInvalidDriver
	}
	if strings.TrimSpace(d.DocumentNumber) == "" {
		return ErrInvalidDriver
	}
	if strings.TrimSpace(d.Phone) == "" {
		return ErrInvalidDriver
	}
	if d.Email != "" && strings.TrimSpace(d.Email) == "" {
		return ErrInvalidDriver
	}
	if strings.TrimSpace(d.LicenseNumber) == "" {
		return ErrInvalidDriver
	}
	if !d.Status.Valid() {
		return ErrInvalidDriver
	}
	return nil
}

func (s DriverStatus) Valid() bool {
	switch s {
	case DriverStatusActive, DriverStatusInactive, DriverStatusSuspended:
		return true
	default:
		return false
	}
}
