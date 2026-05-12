package domain

import "context"

type Repository interface {
	Create(ctx context.Context, vehicle *Vehicle) error
	GetByID(ctx context.Context, companyID string, id string) (*Vehicle, error)
	ListByCompany(ctx context.Context, companyID string) ([]Vehicle, error)
	Update(ctx context.Context, vehicle *Vehicle) error
	MarkInactive(ctx context.Context, companyID string, id string) error
}
