package http

import "time"

type loginRequest struct {
	CompanyID string `json:"company_id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type userResponse struct {
	ID        string    `json:"id"`
	CompanyID string    `json:"company_id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type loginResponse struct {
	AccessToken string       `json:"access_token"`
	TokenType   string       `json:"token_type"`
	User        userResponse `json:"user"`
}

type meResponse struct {
	User userResponse `json:"user"`
}

type errorResponse struct {
	Error string `json:"error"`
}
