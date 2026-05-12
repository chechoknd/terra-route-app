package dto

import "time"

type LoginRequest struct {
	CompanyID string
	Email     string
	Password  string
}

type AuthenticatedUser struct {
	ID        string
	CompanyID string
	Email     string
	FullName  string
	Role      string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type LoginResponse struct {
	User        AuthenticatedUser
	AccessToken string
	TokenType   string
}
