package domain

import "errors"

var (
	ErrInvalidRouteStop       = errors.New("invalid route stop")
	ErrRouteStopNotFound      = errors.New("route stop not found")
	ErrRouteStopAlreadyExists = errors.New("route stop already exists")
)
