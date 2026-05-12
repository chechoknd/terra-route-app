package domain

import "context"

type Repository interface {
	Create(ctx context.Context, companyID string, stop *RouteStop) error
	GetByID(ctx context.Context, companyID string, routeID string, id string) (*RouteStop, error)
	ListByRoute(ctx context.Context, companyID string, routeID string) ([]RouteStop, error)
	Update(ctx context.Context, companyID string, stop *RouteStop) error
	Delete(ctx context.Context, companyID string, routeID string, id string) error
	Reorder(ctx context.Context, companyID string, routeID string, orderedIDs []string) error
}
