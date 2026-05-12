package domain

import "errors"

var (
	ErrInvalidDriver       = errors.New("invalid driver")
	ErrDriverNotFound      = errors.New("driver not found")
	ErrDriverAlreadyExists = errors.New("driver already exists")
)
