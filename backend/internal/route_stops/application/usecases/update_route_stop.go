package usecases

import (
	"context"

	"github.com/terraroute/terra-route/backend/internal/route_stops/application/dto"
	"github.com/terraroute/terra-route/backend/internal/route_stops/domain"
)

type UpdateRouteStopUseCase struct {
	stops domain.Repository
}

func NewUpdateRouteStopUseCase(stops domain.Repository) *UpdateRouteStopUseCase {
	return &UpdateRouteStopUseCase{stops: stops}
}

func (uc *UpdateRouteStopUseCase) Execute(ctx context.Context, req dto.UpdateRouteStopRequest) (*dto.RouteStop, error) {
	if req.CompanyID == "" || req.ID == "" {
		return nil, domain.ErrInvalidRouteStop
	}
	current, err := uc.stops.GetByID(ctx, req.CompanyID, req.RouteID, req.ID)
	if err != nil {
		return nil, err
	}
	stop := &domain.RouteStop{
		ID:        current.ID,
		RouteID:   current.RouteID,
		Name:      req.Name,
		City:      req.City,
		StopOrder: req.StopOrder,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		CreatedAt: current.CreatedAt,
		UpdatedAt: current.UpdatedAt,
	}
	if err := stop.Validate(); err != nil {
		return nil, err
	}
	if err := uc.stops.Update(ctx, req.CompanyID, stop); err != nil {
		return nil, err
	}
	res := routeStopToDTO(*stop)
	return &res, nil
}
