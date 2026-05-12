package usecases

import (
	"context"

	"github.com/terraroute/terra-route/backend/internal/route_stops/application/dto"
	"github.com/terraroute/terra-route/backend/internal/route_stops/domain"
)

type ListRouteStopsUseCase struct {
	stops domain.Repository
}

func NewListRouteStopsUseCase(stops domain.Repository) *ListRouteStopsUseCase {
	return &ListRouteStopsUseCase{stops: stops}
}

func (uc *ListRouteStopsUseCase) Execute(ctx context.Context, req dto.ListRouteStopsRequest) ([]dto.RouteStop, error) {
	if req.CompanyID == "" || req.RouteID == "" {
		return nil, domain.ErrInvalidRouteStop
	}
	stops, err := uc.stops.ListByRoute(ctx, req.CompanyID, req.RouteID)
	if err != nil {
		return nil, err
	}
	res := make([]dto.RouteStop, 0, len(stops))
	for _, stop := range stops {
		res = append(res, routeStopToDTO(stop))
	}
	return res, nil
}
