package usecases

import (
	"context"

	"github.com/terraroute/terra-route/backend/internal/route_stops/application/dto"
	"github.com/terraroute/terra-route/backend/internal/route_stops/domain"
)

type AddRouteStopUseCase struct {
	stops domain.Repository
}

func NewAddRouteStopUseCase(stops domain.Repository) *AddRouteStopUseCase {
	return &AddRouteStopUseCase{stops: stops}
}

func (uc *AddRouteStopUseCase) Execute(ctx context.Context, req dto.AddRouteStopRequest) (*dto.RouteStop, error) {
	if req.CompanyID == "" {
		return nil, domain.ErrInvalidRouteStop
	}
	stop := &domain.RouteStop{
		RouteID:   req.RouteID,
		Name:      req.Name,
		City:      req.City,
		StopOrder: req.StopOrder,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	}
	if err := stop.Validate(); err != nil {
		return nil, err
	}
	if err := uc.stops.Create(ctx, req.CompanyID, stop); err != nil {
		return nil, err
	}
	res := routeStopToDTO(*stop)
	return &res, nil
}
