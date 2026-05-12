package usecases

import (
	"context"

	"github.com/terraroute/terra-route/backend/internal/routes/application/dto"
	"github.com/terraroute/terra-route/backend/internal/routes/domain"
)

type UpdateRouteUseCase struct {
	routes domain.Repository
}

func NewUpdateRouteUseCase(routes domain.Repository) *UpdateRouteUseCase {
	return &UpdateRouteUseCase{routes: routes}
}

func (uc *UpdateRouteUseCase) Execute(ctx context.Context, req dto.UpdateRouteRequest) (*dto.Route, error) {
	if req.ID == "" {
		return nil, domain.ErrInvalidRoute
	}

	current, err := uc.routes.GetByID(ctx, req.CompanyID, req.ID)
	if err != nil {
		return nil, err
	}

	route := &domain.Route{
		ID:                       current.ID,
		CompanyID:                current.CompanyID,
		Name:                     req.Name,
		OriginCity:               req.OriginCity,
		DestinationCity:          req.DestinationCity,
		EstimatedDistanceKM:      req.EstimatedDistanceKM,
		EstimatedDurationMinutes: req.EstimatedDurationMinutes,
		BasePrice:                req.BasePrice,
		Status:                   domain.RouteStatus(req.Status),
		CreatedAt:                current.CreatedAt,
		UpdatedAt:                current.UpdatedAt,
	}
	if err := validateRouteStatusTransition(current.Status, route.Status); err != nil {
		return nil, err
	}
	if err := route.Validate(); err != nil {
		return nil, err
	}

	if err := uc.routes.Update(ctx, route); err != nil {
		return nil, err
	}

	res := routeToDTO(*route)
	return &res, nil
}

func validateRouteStatusTransition(from domain.RouteStatus, to domain.RouteStatus) error {
	if !to.Valid() {
		return domain.ErrInvalidRoute
	}
	if from == domain.RouteStatusArchived && to != domain.RouteStatusArchived {
		return domain.ErrInvalidRoute
	}
	return nil
}
