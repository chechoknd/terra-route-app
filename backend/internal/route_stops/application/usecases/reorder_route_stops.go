package usecases

import (
	"context"

	"github.com/terraroute/terra-route/backend/internal/route_stops/application/dto"
	"github.com/terraroute/terra-route/backend/internal/route_stops/domain"
)

type ReorderRouteStopsUseCase struct {
	stops domain.Repository
}

func NewReorderRouteStopsUseCase(stops domain.Repository) *ReorderRouteStopsUseCase {
	return &ReorderRouteStopsUseCase{stops: stops}
}

func (uc *ReorderRouteStopsUseCase) Execute(ctx context.Context, req dto.ReorderRouteStopsRequest) error {
	if req.CompanyID == "" || req.RouteID == "" || len(req.OrderedIDs) == 0 {
		return domain.ErrInvalidRouteStop
	}
	for _, id := range req.OrderedIDs {
		if id == "" {
			return domain.ErrInvalidRouteStop
		}
	}
	return uc.stops.Reorder(ctx, req.CompanyID, req.RouteID, req.OrderedIDs)
}
