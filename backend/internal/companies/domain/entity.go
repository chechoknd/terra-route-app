package domain

import (
	"strings"
	"time"
)

const (
	CompanyStatusActive    = "active"
	CompanyStatusInactive  = "inactive"
	CompanyStatusSuspended = "suspended"
)

type Company struct {
	ID        string
	Name      string
	Slug      string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (c Company) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return ErrInvalidCompany
	}
	if strings.TrimSpace(c.Slug) == "" {
		return ErrInvalidCompany
	}
	if !validCompanyStatus(c.Status) {
		return ErrInvalidCompany
	}
	return nil
}

func validCompanyStatus(status string) bool {
	switch status {
	case CompanyStatusActive, CompanyStatusInactive, CompanyStatusSuspended:
		return true
	default:
		return false
	}
}
