package usecases

import (
	"context"

	"github.com/terraroute/terra-route/backend/internal/routes/application/dto"
	"github.com/terraroute/terra-route/backend/internal/routes/domain"
)

type CreateRouteUseCase struct {
	routes domain.Repository
}

func NewCreateRouteUseCase(routes domain.Repository) *CreateRouteUseCase {
	return &CreateRouteUseCase{routes: routes}
}

func (uc *CreateRouteUseCase) Execute(ctx context.Context, req dto.CreateRouteRequest) (*dto.Route, error) {
	route := &domain.Route{
		CompanyID:                req.CompanyID,
		Name:                     req.Name,
		OriginCity:               req.OriginCity,
		DestinationCity:          req.DestinationCity,
		EstimatedDistanceKM:      req.EstimatedDistanceKM,
		EstimatedDurationMinutes: req.EstimatedDurationMinutes,
		BasePrice:                req.BasePrice,
		Status:                   routeStatusOrDefault(req.Status),
	}
	if err := route.Validate(); err != nil {
		return nil, err
	}

	if err := uc.routes.Create(ctx, route); err != nil {
		return nil, err
	}

	res := routeToDTO(*route)
	return &res, nil
}
