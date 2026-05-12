package domain

import "errors"

var (
	ErrInvalidRoute       = errors.New("invalid route")
	ErrRouteNotFound      = errors.New("route not found")
	ErrRouteAlreadyExists = errors.New("route already exists")
)
