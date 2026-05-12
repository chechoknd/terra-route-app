package domain

import "errors"

var (
	ErrInvalidVehicle       = errors.New("invalid vehicle")
	ErrVehicleNotFound      = errors.New("vehicle not found")
	ErrVehicleAlreadyExists = errors.New("vehicle already exists")
)
