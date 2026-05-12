package domain

import "context"

type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, companyID string, id string) (*User, error)
	GetByEmail(ctx context.Context, companyID string, email string) (*User, error)
	ListByCompany(ctx context.Context, companyID string) ([]User, error)
	Update(ctx context.Context, user *User) error
}
