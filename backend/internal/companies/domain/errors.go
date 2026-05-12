package domain

import "errors"

var (
	ErrInvalidCompany  = errors.New("invalid company")
	ErrCompanyNotFound = errors.New("company not found")
)
