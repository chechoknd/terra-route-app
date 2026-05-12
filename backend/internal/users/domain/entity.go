package domain

import (
	"strings"
	"time"
)

const (
	UserRoleSuperAdmin   = "super_admin"
	UserRoleCompanyAdmin = "company_admin"
	UserRoleOperator     = "operator"
	UserRoleDriver       = "driver"

	UserStatusActive    = "active"
	UserStatusInactive  = "inactive"
	UserStatusSuspended = "suspended"
)

type User struct {
	ID           string
	CompanyID    string
	Email        string
	FullName     string
	Role         string
	Status       string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (u User) Validate() error {
	if strings.TrimSpace(u.Email) == "" {
		return ErrInvalidUser
	}
	if strings.TrimSpace(u.FullName) == "" {
		return ErrInvalidUser
	}
	if strings.TrimSpace(u.PasswordHash) == "" {
		return ErrInvalidUser
	}
	if !validUserRole(u.Role) {
		return ErrInvalidUser
	}
	if !validUserStatus(u.Status) {
		return ErrInvalidUser
	}
	if u.Role == UserRoleSuperAdmin {
		if strings.TrimSpace(u.CompanyID) != "" {
			return ErrInvalidUser
		}
		return nil
	}
	if strings.TrimSpace(u.CompanyID) == "" {
		return ErrInvalidUser
	}
	return nil
}

func validUserRole(role string) bool {
	switch role {
	case UserRoleSuperAdmin, UserRoleCompanyAdmin, UserRoleOperator, UserRoleDriver:
		return true
	default:
		return false
	}
}

func validUserStatus(status string) bool {
	switch status {
	case UserStatusActive, UserStatusInactive, UserStatusSuspended:
		return true
	default:
		return false
	}
}
