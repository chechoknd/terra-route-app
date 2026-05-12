package domain

import (
	"context"
	"time"
)

type TokenSubject struct {
	UserID    string
	CompanyID string
	Role      string
}

type TokenClaims struct {
	UserID    string
	CompanyID string
	Role      string
	ExpiresAt time.Time
}

type TokenService interface {
	Generate(ctx context.Context, subject TokenSubject) (string, error)
	Validate(ctx context.Context, token string) (*TokenClaims, error)
}
