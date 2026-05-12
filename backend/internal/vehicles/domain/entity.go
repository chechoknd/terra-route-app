package domain

import (
	"strings"
	"time"
)

type VehicleStatus string

const (
	VehicleStatusActive      VehicleStatus = "active"
	VehicleStatusInactive    VehicleStatus = "inactive"
	VehicleStatusMaintenance VehicleStatus = "maintenance"
	VehicleStatusUnavailable VehicleStatus = "unavailable"
)

type Vehicle struct {
	ID           string
	CompanyID    string
	Plate        string
	InternalCode string
	VehicleType  string
	Brand        string
	Model        string
	Capacity     int
	Status       VehicleStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (v Vehicle) Validate() error {
	if strings.TrimSpace(v.CompanyID) == "" {
		return ErrInvalidVehicle
	}
	if strings.TrimSpace(v.Plate) == "" {
		return ErrInvalidVehicle
	}
	if strings.TrimSpace(v.InternalCode) == "" {
		return ErrInvalidVehicle
	}
	if strings.TrimSpace(v.VehicleType) == "" {
		return ErrInvalidVehicle
	}
	if strings.TrimSpace(v.Brand) == "" {
		return ErrInvalidVehicle
	}
	if strings.TrimSpace(v.Model) == "" {
		return ErrInvalidVehicle
	}
	if v.Capacity <= 0 {
		return ErrInvalidVehicle
	}
	if !v.Status.Valid() {
		return ErrInvalidVehicle
	}
	return nil
}

func (s VehicleStatus) Valid() bool {
	switch s {
	case VehicleStatusActive, VehicleStatusInactive, VehicleStatusMaintenance, VehicleStatusUnavailable:
		return true
	default:
		return false
	}
}
