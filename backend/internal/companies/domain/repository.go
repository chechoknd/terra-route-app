package domain

import "context"

type Repository interface {
	Create(ctx context.Context, company *Company) error
	GetByID(ctx context.Context, id string) (*Company, error)
	GetBySlug(ctx context.Context, slug string) (*Company, error)
	Update(ctx context.Context, company *Company) error
}
