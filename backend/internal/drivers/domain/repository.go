package domain

import "context"

type Repository interface {
	Create(ctx context.Context, driver *Driver) error
	GetByID(ctx context.Context, companyID string, id string) (*Driver, error)
	ListByCompany(ctx context.Context, companyID string) ([]Driver, error)
	Update(ctx context.Context, driver *Driver) error
	Deactivate(ctx context.Context, companyID string, id string) error
}
