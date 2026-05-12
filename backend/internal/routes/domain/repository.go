package domain

import "context"

type Repository interface {
	Create(ctx context.Context, route *Route) error
	GetByID(ctx context.Context, companyID string, id string) (*Route, error)
	ListByCompany(ctx context.Context, companyID string) ([]Route, error)
	Update(ctx context.Context, route *Route) error
	Archive(ctx context.Context, companyID string, id string) error
}
